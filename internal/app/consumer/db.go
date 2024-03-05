package consumer

import (
	"github.com/arslanovdi/logistic-package-api/internal/app/repo"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"sync"
	"time"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n     uint64
	event chan<- model.PackageEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64                    // кол-во потоков (горутин)
	events    chan<- model.PackageEvent // канал событий
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumerTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.PackageEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

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
					events, err := c.repo.Lock(c.batchSize) // берем события из базы
					if err != nil {
						continue
					}
					for _, event := range events { // передаем события в канал
						c.event <- event
					}
				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
