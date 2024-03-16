package main

import (
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
