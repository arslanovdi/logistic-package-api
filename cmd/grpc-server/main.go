package main

import (
	"flag"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

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

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Warn("Failed init configuration", slog.Any("error", err))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

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

	*migration = true // todo: need to delete this line for homework-4
	if *migration {
		//goose.SetBaseFS(migrations.EmbedFS) // можно примонтировать как файлы с миграциями
		//if err = goose.Up(db.DB, "."); err != nil {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Error("Migration failed", slog.Any("error", err))
			return
		}

	}

	//repo := repo.NewRepo(db, batchSize)

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		log.Error("Failed init tracing", slog.Any("error", err))

		return
	}
	defer tracing.Close()

	if err := server.NewGrpcServer(db, batchSize).Start(&cfg); err != nil {
		log.Error("Failed creating gRPC server", slog.Any("error", err))

		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-stop:
		slog.Info("Graceful shutdown")
		goose.Down(db.DB, cfg.Database.Migrations)
		slog.Info("Application stopped")
		return
	}
}
