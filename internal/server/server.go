package server

import (
	"context"
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/arslanovdi/logistic-package-api/internal/api"
	"github.com/arslanovdi/logistic-package-api/internal/app/repo"
	"github.com/arslanovdi/logistic-package-api/internal/config"
	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

// GrpcServer is gRPC server
type GrpcServer struct {
	db        *sqlx.DB
	batchSize uint
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(db *sqlx.DB, batchSize uint) *GrpcServer {
	return &GrpcServer{
		db:        db,
		batchSize: batchSize,
	}
}

// Start method runs server
func (s *GrpcServer) Start(cfg *config.Config) error {

	log := slog.With("func", "server.Start")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	gatewayServer := createGatewayServer(grpcAddr, gatewayAddr)

	go func() {
		log.Info("Gateway server is running", slog.String("address", gatewayAddr))
		log.Info("Swagger server is running", slog.String("address", gatewayAddr+"/swagger-ui/"))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running gateway server", slog.Any("error", err))
			cancel()
		}
	}()

	metricsServer := createMetricsServer(cfg)

	go func() {
		log.Info("Metrics server is running", slog.String("address", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running metrics server", slog.Any("error", err))
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := createStatusServer(cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)
		log.Info("Status server is running", slog.String("address", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running status server", slog.Any("error", err))
		}
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer( // последовательное исполнение middleware, с общим контекстом
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	r := repo.NewRepo(s.db, s.batchSize)

	pb.RegisterLogisticPackageApiServiceServer(grpcServer, api.NewPackageAPI(r)) // регистрируем имплементацию интерфейса в gRPC-сервере
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		log.Info("GRPC Server is listening", slog.String("address", grpcAddr))
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("Failed running gRPC server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Info("The service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer) // в дебаге регестрируем отражение методов gRPC-сервера: предоставляет сведения о публично доступных методах
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info("Gracefully shutdown", slog.Any("signal.Notify", v))
	case done := <-ctx.Done():
		log.Warn("ctx.Done", slog.Any("error", done))
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		log.Error("gatewayServer.Shutdown", slog.Any("error", err))
	} else {
		log.Info("gatewayServer shut down correctly")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		log.Error("statusServer.Shutdown", slog.Any("error", err))
	} else {
		log.Info("statusServer shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Error("metricsServer.Shutdown", slog.Any("error", err))
	} else {
		log.Info("metricsServer shut down correctly")
	}

	grpcServer.GracefulStop()
	log.Info("grpcServer shut down correctly")

	return nil
}
