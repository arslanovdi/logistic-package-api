package repo

import "github.com/arslanovdi/logistic-package-api/internal/model"

type EventRepo interface {
	// Lock заблокировать в БД n записей
	Lock(n uint64) ([]model.PackageEvent, error) // TODO в имплементации реализовать загрузку зависших залоченных событий, с отсечкой по времени
	// Unlock разблокировать в БД n записей
	Unlock(eventID []uint64) error

	// Add добавить в БД n записей
	Add(event []model.PackageEvent) error
	// Remove удалить из БД n записей
	Remove(eventIDs []uint64) error
}
