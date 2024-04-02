package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/database/postgres"
	"github.com/arslanovdi/logistic-package-api/internal/logger"
	"github.com/arslanovdi/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package-api/internal/tracer"
	"github.com/jackc/pgx/v5/stdlib"
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

	if err1 := config.ReadConfigYML("config.yml"); err1 != nil {
		log.Warn("Failed init configuration", slog.String("error", err1.Error()))
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

	dbpool := database.MustGetPgxPool(context.Background())
	defer dbpool.Close()

	repo := postgres.NewPostgresRepo(dbpool, batchSize)
	packageService := service.NewPackageService(dbpool, repo)

	migration := flag.Bool("migration", true, "Defines the migration start option") // миграцию запускаем параметром из командной строки -migration
	flag.Parse()

	if *migration {
		log.Info("Migration started")
		if err := goose.Up(stdlib.OpenDBFromPool(dbpool), // получаем соединение с базой данных из пула
			cfg.Database.Migrations); err != nil {
			log.Warn("Migration failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		//fakedata.Generate(100, repo)
	}

	ctxTrace, cancelTrace := context.WithCancel(context.Background())
	defer cancelTrace()
	trace, err := tracer.NewTracer(ctxTrace)
	if err != nil {
		log.Warn("Failed to init tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}

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
	select {
	case <-ctxServer.Done(): // ошибка при старте любого из серверов
		log.Warn("Fail on start servers")
	case <-stop:
		slog.Info("Graceful shutdown")
		isReady.Store(false)

		goose.Down(stdlib.OpenDBFromPool(dbpool), cfg.Database.Migrations)

		if err := grpcServer.Stop(); err != nil {
			log.Error("Failed to stop gRPC server", slog.String("error", err.Error()))
		}
		metricsServer.Stop(ctxServer)
		statusServer.Stop(ctxServer)
		gatewayServer.Stop(ctxServer)
		if err := trace.Shutdown(ctxTrace); err != nil {
			log.Error("Error shutting down tracer provider", slog.String("error", err.Error()))
		}

		slog.Info("Application stopped")
	}
}
