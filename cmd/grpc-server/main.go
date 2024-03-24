package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/logger"
	"github.com/arslanovdi/logistic-package-api/internal/service"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/arslanovdi/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package-api/internal/database"
	"github.com/arslanovdi/logistic-package-api/internal/server"
	"github.com/arslanovdi/logistic-package-api/internal/tracer"
)

var (
	batchSize uint = 2
)

func main() {
	logger.InitializeLogger()
	log := slog.With("func", "grpc-server.main")

	startCtx, cancel := context.WithTimeout(context.Background(), time.Minute) // контекст запуска приложения
	defer cancel()
	go func() {
		<-startCtx.Done()
		if errors.Is(startCtx.Err(), context.DeadlineExceeded) { // приложение зависло при запуске
			log.Warn("Application startup time exceeded")
			os.Exit(1)
		}
	}()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Warn("Failed init configuration", slog.Any("error", err))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

	if cfg.Project.Debug {
		logger.SetLogLevel(slog.LevelDebug)
	} else {
		logger.SetLogLevel(slog.LevelInfo)
	}

	log.Info(fmt.Sprintf("Starting service %s", cfg.Project.Name),
		slog.String("version", cfg.Project.Version),
		slog.String("commitHash", cfg.Project.CommitHash),
		slog.Bool("debug", cfg.Project.Debug),
		slog.String("environment", cfg.Project.Environment),
	)

	db, err := database.NewPostgres()
	if err != nil {
		log.Warn("Failed init postgres", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	migration := flag.Bool("migration", true, "Defines the migration start option") // миграцию запускаем параметром из командной строки -migration
	flag.Parse()

	if *migration {
		log.Info("Migration started")
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Warn("Migration failed", slog.Any("error", err))
			os.Exit(1)
		}
	}

	repo := database.NewRepo(db, batchSize)
	packageService := service.NewPackageService(db, repo)

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		log.Error("Failed init tracing", slog.Any("error", err))
		os.Exit(1)
	}
	defer tracing.Close()

	ctxServer, cancelServer := context.WithCancel(context.Background()) // контекст запуска серверов, при ошибке в любом из серверов контекст отменяется
	isReady := &atomic.Value{}
	isReady.Store(false)

	go func() { // TODO отсечка статус сервера
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Info("The service is ready to accept requests")
	}()

	grpcServer := server.NewGrpcServer(packageService, batchSize)
	grpcServer.Start(cancelServer)

	metricsServer := server.NewMetricsServer()
	metricsServer.Start(cancelServer)

	statusServer := server.NewStatusServer(isReady)
	statusServer.Start(cancelServer)

	gatewayServer := server.NewGatewayServer()
	gatewayServer.Start(cancelServer)

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-ctxServer.Done(): // ошибка при старте любого из серверов
			log.Warn("Fail on start servers")
			os.Exit(1)
		case <-stop:
			slog.Info("Graceful shutdown")
			isReady.Store(false)
			//goose.Down(db.DB, cfg.Database.Migrations)
			grpcServer.Stop()
			metricsServer.Stop(ctxServer)
			statusServer.Stop(ctxServer)
			gatewayServer.Stop(ctxServer)
			slog.Info("Application stopped")
			return
		}
	}
}
