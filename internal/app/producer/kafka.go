package producer

import (
	"github.com/arslanovdi/logistic-package-api/internal/app/repo"
	"github.com/arslanovdi/logistic-package-api/internal/app/sender"
	"github.com/arslanovdi/logistic-package-api/internal/app/workerpool"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"log/slog"
	"sync"
	"time"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.PackageEvent

	repo       repo.EventRepo
	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
}

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
					if err := p.sender.Send(&event); err != nil {
						p.workerPool.Submit(func() { // снимаем блокировку с события в БД, т.к. отправка в кавку неудачная
							err := p.repo.Unlock([]uint64{event.ID})
							if err != nil {
								slog.Error("Ошибка при снятии блокировки с события в БД", err)
							}
						})
					} else {
						p.workerPool.Submit(func() { // удаляем событие из БД, т.к. оно обработано и отправлено в кавку
							err := p.repo.Remove([]uint64{event.ID})
							if err != nil {
								slog.Error("Ошибка при удалении события из БД", err)
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

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
	slog.Debug("Kafka producer stopped")
}
