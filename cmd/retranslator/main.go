// Сервис для пересылки событий из базы данных в кафку (outbox pattern)
package main

import (
	"context"
	"errors"
	"github.com/arslanovdi/logistic-package-api/internal/database"
	"github.com/arslanovdi/logistic-package-api/internal/database/postgres"
	"github.com/arslanovdi/logistic-package-api/internal/logger"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/retranslator"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const level = slog.LevelDebug // log level
//const batchsize = 10

func main() {
	logger.SetLogLevel(level)
	log := slog.With("func", "kafkaRetranslator.main")

	startCtx, cancel := context.WithTimeout(context.Background(), time.Minute) // контекст запуска приложения
	defer cancel()
	go func() {
		<-startCtx.Done()
		if errors.Is(startCtx.Err(), context.DeadlineExceeded) { // приложение зависло при запуске
			log.Warn("Application startup time exceeded")
			os.Exit(1)
		}
	}()

	pool := database.MustGetPgxPool(startCtx)

	cfg := retranslator.Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  28,
		WorkerCount:    2,
		Repo:           postgres.NewPostgresRepo(pool /*, batchsize*/),
		Sender:         nil,
	}

	kafkaRetranslator := retranslator.NewRetranslator(cfg)
	kafkaRetranslator.Start()
	slog.Info("Retranslator started")

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	slog.Info("Graceful shutdown")
	kafkaRetranslator.Stop()
	pool.Close()
	slog.Info("Application stopped")
}
