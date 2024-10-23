// Package consumer принимает события из базы данных и отправляет их в канал
package consumer

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package-api/internal/server"
	"log/slog"
	"sync"
	"time"
)

// Consumer читает из базы данных события в n потоков и отправляет их в канал
type Consumer interface {
	Start()
	Stop()
}

type consumer struct {
	repo      repo.EventRepo
	n         uint64
	batchSize uint64
	done      chan bool
	event     chan<- model.PackageEvent
	timeout   time.Duration
	wg        *sync.WaitGroup
}

// NewDbConsumer конструктор
func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumerTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.PackageEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	slog.Debug("db consumer created")
	return &consumer{
		n:         n,
		event:     events,
		repo:      repo,
		batchSize: batchSize,
		timeout:   consumerTimeout,
		wg:        wg,
		done:      done,
	}
}

// Start starts consumer
func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ { // запускаем n горутин
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout) // тикер
			for {
				select {
				case <-ticker.C: // если тикер сработал

					events, err := c.repo.Lock(context.TODO(), c.batchSize) // берем события из базы
					if err != nil {
						slog.Error("Error getting events from db", slog.String("error", err.Error()))
						continue
					}
					server.RetranslatorEvents.Add(float64(len(events))) // метрика, кол-во обрабатываемых событий, прибавляем к счетчику
					for _, event := range events {                      // передаем события в канал
						c.event <- event
					}
				case <-c.done:
					return
				}
			}
		}()
	}
}

// Stop consumer
func (c *consumer) Stop() {
	close(c.done)
	c.wg.Wait()
	slog.Debug("db consumer stopped")
}
