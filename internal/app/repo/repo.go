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
	Lock(n uint64) ([]model.PackageEvent, error)
	Unlock(eventID []uint64) error

	Add(event []model.PackageEvent) error
	Remove(eventIDs []uint64) error
>>>>>>> 7177f39 (Initial commit)
}
