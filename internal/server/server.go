// Package server - все http сервисы приложения
package server

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/service"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"time"

	"github.com/arslanovdi/logistic-package-api/internal/api"
	"github.com/arslanovdi/logistic-package-api/internal/config"
	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// GrpcServer is gRPC server
type GrpcServer struct {
	server    *grpc.Server
	lis       net.Listener
	batchSize uint
}

// grpcMiddleware Перехватчик унарных методов, считаем метрики
func grpcMiddleware(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	d := time.Now()

	m, err := handler(ctx, req)

	GRPC2.Observe(time.Since(d).Seconds())

	if status.Code(err) == codes.NotFound {
		GRPCNotFoundCounter.Inc()
	}

	CRUDCounter.Inc()

	return m, err
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(packageService *service.PackageService, batchSize uint) *GrpcServer {

	cfg := config.GetConfigInstance()

	s := &GrpcServer{
		batchSize: batchSize,
	}

	s.server = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // openTelemetry трассировка
		grpc.ChainUnaryInterceptor( // последовательное исполнение middleware, с общим контекстом
			grpcMiddleware,
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	pb.RegisterLogisticPackageApiServiceServer(s.server, api.NewPackageAPI(packageService)) // регистрируем имплементацию интерфейса в gRPC-сервере
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s.server)

	if cfg.Project.Debug {
		reflection.Register(s.server) // в дебаге регестрируем отражение методов gRPC-сервера: предоставляет сведения о публично доступных методах
	}

	return s
}

// Start method runs server
func (s *GrpcServer) Start(cancelFunc context.CancelFunc) {

	log := slog.With("func", "GrpcServer.Start")

	cfg := config.GetConfigInstance()

	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)

	var err1 error
	s.lis, err1 = net.Listen("tcp", grpcAddr)
	if err1 != nil {
		log.Error("failed to listen", slog.String("error", err1.Error()))
		cancelFunc()
	}

	go func() {
		log.Info("GRPC Server is listening", slog.String("address", grpcAddr))
		if err2 := s.server.Serve(s.lis); err2 != nil {
			log.Error("Failed running gRPC server", slog.String("error", err2.Error()))
			cancelFunc()
		}
	}()
}

// Stop - stop gRPC server
func (s *GrpcServer) Stop() error {
	s.server.GracefulStop()
	err := s.lis.Close()
	if err != nil {
		return err
	}
	return nil
}
