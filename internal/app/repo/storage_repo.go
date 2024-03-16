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
func NewRepo(db *sqlx.DB, batchSize uint) *repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) Describe(ctx context.Context, templateID uint64) (*model.Package, error) {
	return nil, nil
}
