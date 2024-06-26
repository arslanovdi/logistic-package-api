package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"net/http"

	"github.com/arslanovdi/logistic-package-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var GRPCCounter = promauto.NewCounter(prometheus.CounterOpts{ // Инициализируем счётчик gRPC вызовов в prometheus
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "grpc_total",
	Help:      "Total gRPC calls",
})

var GRPC2 = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "logistic",
	Subsystem:  "package",
	Name:       "grpc2",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 0,5 медиана, 0,9 90%-й квантиль, 0,99 99%-й квантиль
}, []string{"method"}, // метка для метрики, для каждой метки будет свой график
)

type metricsServer struct {
	server *http.Server
}

func NewMetricsServer() *metricsServer {

	cfg := config.GetConfigInstance()

	addr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metrics := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &metricsServer{
		server: metrics,
	}
}

func (s *metricsServer) Start(cancelFunc context.CancelFunc) {

	log := slog.With("func", "MetricsServer.Start")

	cfg := config.GetConfigInstance()

	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	go func() {
		log.Info("Metrics server is running", slog.String("address", metricsAddr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running metrics server", slog.String("error", err.Error()))
			cancelFunc()
		}
	}()
}

func (s *metricsServer) Stop(ctx context.Context) {

	log := slog.With("func", "MetricsServer.Stop")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("metricsServer.Shutdown", slog.String("error", err.Error()))
	} else {
		log.Info("metricsServer shut down correctly")
	}
}
