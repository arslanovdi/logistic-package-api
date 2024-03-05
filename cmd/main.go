package main

import (
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/app/retranslator"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		retranslator.Close()
		fmt.Println("Graceful shutdown")
		return
	}
}
