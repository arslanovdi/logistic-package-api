// Package retranslator get events from database (consumer) and send to kafka (producer)
package retranslator

import (
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/consumer"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/producer"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package-api/internal/outbox/workerpool"
	"log/slog"
	"time"
)

// Retranslator считывает события из БД и отправляет в кафку
type Retranslator interface {
	Start()
	Stop()
}

// Config конфигурация Retranslator
type Config struct {
	ChannelSize uint64

	ConsumerCount  uint64
	ConsumeSize    uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount   uint64

	Repo   repo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events     chan model.PackageEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

// NewRetranslator конструктор
func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.PackageEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	dbconsumer := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)
	kafkaproducer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		cfg.Repo,
		workerPool)

	slog.Debug("Retranslator created")
	return &retranslator{
		events:     events,
		consumer:   dbconsumer,
		producer:   kafkaproducer,
		workerPool: workerPool,
	}
}

// Start запускает пул воркеров, считывает события из БД и отправляет в кафку
func (r *retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

// Stop останавливает пул воркеров
func (r *retranslator) Stop() {
	r.consumer.Stop()
	r.producer.Stop()
	r.workerPool.StopWait()

	slog.Debug("Retranslator stopped")
}
