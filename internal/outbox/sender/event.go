// Package sender send events to kafka
package sender

import "github.com/arslanovdi/logistic-package-api/internal/model"

// EventSender - интерфейс для отправки событий в кафку
type EventSender interface {
	// Send отправить событие в кавку
	Send(pkg *model.PackageEvent) error
}
