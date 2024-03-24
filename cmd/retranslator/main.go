package main

import (
	"context"
	"errors"
	"github.com/arslanovdi/logistic-package-api/internal/app/retranslator"
	"github.com/arslanovdi/logistic-package-api/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const level = slog.LevelDebug // log level

func main() {
	logger.SetLogLevel(level)
	log := slog.With("func", "retranslator.main")

	startCtx, cancel := context.WithTimeout(context.Background(), time.Minute) // контекст запуска приложения
	defer cancel()
	go func() {
		<-startCtx.Done()
		if errors.Is(startCtx.Err(), context.DeadlineExceeded) { // приложение зависло при запуске
			log.Warn("Application startup time exceeded")
			os.Exit(1)
		}
	}()

	cfg := retranslator.Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  28,
		WorkerCount:    2,
		Repo:           nil,
		Sender:         nil,
	}

	retranslator := retranslator.NewRetranslator(cfg)
	retranslator.Start()
	slog.Info("Retranslator started")

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-stop:
		slog.Info("Graceful shutdown")
		retranslator.Close()
		slog.Info("Application stopped")
		return
	}
}
