package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/arslanovdi/logistic-package-api/internal/config"
)

type statusServer struct {
	server *http.Server
}

func NewStatusServer(isReady *atomic.Value) *statusServer {

	cfg := config.GetConfigInstance()

	statusAddr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	mux := http.DefaultServeMux

	mux.HandleFunc(cfg.Status.LivenessPath, livenessHandler)
	mux.HandleFunc(cfg.Status.ReadinessPath, readinessHandler(isReady))
	mux.HandleFunc(cfg.Status.VersionPath, versionHandler(&cfg))

	server := &http.Server{
		Addr:    statusAddr,
		Handler: mux,
	}

	return &statusServer{
		server: server,
	}
}

func (s *statusServer) Start(cancelFunc context.CancelFunc) {

	log := slog.With("func", "StatusServer.Start")

	cfg := config.GetConfigInstance()

	statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	go func() {
		log.Info("Status server is running", slog.String("address", statusAdrr))
		if err1 := s.server.ListenAndServe(); err1 != nil && !errors.Is(err1, http.ErrServerClosed) {
			log.Error("Failed running status server", slog.String("error", err1.Error()))
			cancelFunc()
		}
	}()
}

func (s *statusServer) Stop(ctx context.Context) {

	log := slog.With("func", "StatusServer.Stop")

	if err1 := s.server.Shutdown(ctx); err1 != nil {
		log.Error("statusServer.Shutdown", slog.String("error", err1.Error()))
	} else {
		log.Info("statusServer shut down correctly")
	}
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func versionHandler(cfg *config.Config) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {

		log := slog.With("func", "versionHandler")

		data := map[string]interface{}{
			"name":        cfg.Project.Name,
			"debug":       cfg.Project.Debug,
			"environment": cfg.Project.Environment,
			"version":     cfg.Project.Version,
			"commitHash":  cfg.Project.CommitHash,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err1 := json.NewEncoder(w).Encode(data); err1 != nil {
			log.Error("Service information encoding error", slog.String("error", err1.Error()))
		}
	}
}
