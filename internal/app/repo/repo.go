package repo

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Repo is DAO for Template
type Repo interface {
	Describe(ctx context.Context, PackageID uint64) (*model.Package, error) // описание
	/*	Remove(ctx context.Context, packageID uint64) (bool, error)
		Create(ctx context.Context, pkg model.Package) (uint64, error)
		List(ctx context.Context, cursor uint64, limit uint64) ([]model.Package, error)
		Update(ctx context.Context, packageID uint64, pkg model.Package) error
		Get(ctx context.Context, cursor uint64) (model.Package, error)*/
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) Describe(ctx context.Context, templateID uint64) (*model.Package, error) {
	return nil, nil
}

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
