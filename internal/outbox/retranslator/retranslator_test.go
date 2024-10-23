package retranslator

import (
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_retranslator(t *testing.T) {

	repo := mocks.NewEventRepo(t)
	sender := mocks.NewEventSender(t)

	repo.EXPECT().Lock(mock.Anything, mock.AnythingOfType("uint64")).Return([]model.PackageEvent{}, nil) // проверяем, что метод Lock вызывался корректно

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(cfg.ConsumeTimeout + time.Second) // ждем пока тикнет консьюмер, на 1 секунду больше таймера консьюмера
	retranslator.Stop()
}
