package repo

import "github.com/arslanovdi/logistic-package-api/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.PackageEvent, error)
	Unlock(eventID []uint64) error

	Add(event []model.PackageEvent) error
	Remove(eventIDs []uint64) error
}
