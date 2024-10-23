// Package producer provides kafka producer
package producer

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/workerpool"
	"log/slog"
	"sync"
	"time"
)

// Producer читает из канала событий и отправляет в кафку, в n потоков
type Producer interface {
	Start()
	Stop()
}

type producer struct {
	sender     sender.EventSender
	repo       repo.EventRepo
	n          uint64 // кол-во потоков
	workerPool *workerpool.WorkerPool
	wg         *sync.WaitGroup
	done       chan bool
	events     <-chan model.PackageEvent
	timeout    time.Duration
}

// NewKafkaProducer конструктор
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.PackageEvent,
	repo repo.EventRepo,
	workerPool *workerpool.WorkerPool,
) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	slog.Debug("kafka producer created")

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		repo:       repo,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
		timeout:    5 * time.Second,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err1 := p.sender.Send(&event); err1 != nil {
						p.workerPool.Submit(func() { // снимаем блокировку с события в БД, т.к. отправка в кавку неудачная
							err2 := p.repo.Unlock(context.TODO(), []uint64{event.ID})
							if err2 != nil {
								slog.Error("Ошибка при снятии блокировки с события в БД", slog.String("error", err2.Error()))
							}
						})
					} else {
						p.workerPool.Submit(func() { // удаляем событие из БД, т.к. оно обработано и отправлено в кавку
							err3 := p.repo.Remove(context.TODO(), []uint64{event.ID})
							if err3 != nil {
								slog.Error("Ошибка при удалении события из БД", slog.String("error", err3.Error()))
							}
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

// Stop останавливает пул и ждет, пока все задачи будут выполнены
func (p *producer) Stop() {
	close(p.done)
	p.wg.Wait()
	slog.Debug("Kafka producer stopped")
}
