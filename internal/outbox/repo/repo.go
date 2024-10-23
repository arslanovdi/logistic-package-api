// Package repo - работа с событиями в БД
package repo

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/model"
)

// EventRepo - интерфейс работы с БД событий.
type EventRepo interface {
	// Lock заблокировать в БД n записей
	Lock(ctx context.Context, n uint64) ([]model.PackageEvent, error)
	// Unlock разблокировать в БД n записей
	Unlock(ctx context.Context, eventID []uint64) error
	// Remove удалить из БД n записей
	Remove(ctx context.Context, eventIDs []uint64) error
}
