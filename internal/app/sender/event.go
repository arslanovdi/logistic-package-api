package sender

import "github.com/arslanovdi/logistic-package-api/internal/model"

type EventSender interface {
	Send(pkg *model.PackageEvent) error
}
