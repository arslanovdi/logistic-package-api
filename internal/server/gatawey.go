package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/config"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"

	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_requests_total",
		Help: "The total number of incoming HTTP requests",
	})
)

type gatewayServer struct {
	server *http.Server
}

// NewGatewayServer returns HTTP gRPC-gateway server
func NewGatewayServer() *gatewayServer {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.

	cfg := config.GetConfigInstance()
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)

	log := slog.With("func", "server.createGatewayServer")

	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Warn("Failed to dial gRPC server",
			slog.String("func", "createGatewayServer"),
			slog.Any("error", err))
	}

	rmux := runtime.NewServeMux()
	if err := pb.RegisterLogisticPackageApiServiceHandler(context.Background(), rmux, conn); err != nil {
		log.Warn("Failed registration handler",
			slog.String("func", "createGatewayServer"),
			slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) { // Подменяем swagger.json, указанный в файле swagger-initializer.js сгенерированным logistic_package_api.swagger.json
		http.ServeFile(w, r, "./swagger/logistic_package_api.swagger.json")
	})

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/"))))

	server := &http.Server{
		Addr:    gatewayAddr,
		Handler: tracingWrapper(mux), // трэйсы, включая сваггер запросы
	}

	return &gatewayServer{
		server: server,
	}
}

func (s *gatewayServer) Start(cancelFunc context.CancelFunc) {
	log := slog.With("func", "GatewayServer.Start")

	cfg := config.GetConfigInstance()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)

	go func() {
		log.Info("Gateway server is running", slog.String("address", gatewayAddr))
		log.Info("Swagger server is running", slog.String("address", gatewayAddr+"/swagger-ui/"))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running gateway server", slog.Any("error", err))
			cancelFunc()
		}
	}()
}

func (s *gatewayServer) Stop(ctx context.Context) {
	log := slog.With("func", "GatewayServer.Stop")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("gatewayServer.Shutdown", slog.Any("error", err))
	} else {
		log.Info("gatewayServer shut down correctly")
	}
}

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		httpTotalRequests.Inc()
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || errors.Is(err, opentracing.ErrSpanContextNotFound) {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}
