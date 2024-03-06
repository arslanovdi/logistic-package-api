package sender

import "github.com/arslanovdi/logistic-package-api/internal/model"

type EventSender interface {
	// Send отправить событие в кавку
	Send(pkg *model.PackageEvent) error
}
