package repo

import "github.com/arslanovdi/logistic-package-api/internal/model"

type EventRepo interface {
	// Lock заблокировать в БД n записей
	Lock(n uint64) ([]model.PackageEvent, error)
	// Unlock разблокировать в БД n записей
	Unlock(eventID []uint64) error
	// Remove удалить из БД n записей
	Remove(eventIDs []uint64) error
}
