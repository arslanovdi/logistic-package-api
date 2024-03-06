package repo

<<<<<<< HEAD
import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/omp-template-api/internal/model"
)

// Repo is DAO for Template
type Repo interface {
	DescribeTemplate(ctx context.Context, templateID uint64) (*model.Template, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeTemplate(ctx context.Context, templateID uint64) (*model.Template, error) {
	return nil, nil
=======
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
>>>>>>> 7177f39 (Initial commit)
}
