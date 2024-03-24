package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Repo is DAO for Template
type Repo interface {
	Create(ctx context.Context, pkg model.Package) (*uint64, error)
	DeletePackage(ctx context.Context, id uint64) (bool, error)
	GetPackage(ctx context.Context, id uint64) (*model.Package, error)
	ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error)
	UpdatePackage(ctx context.Context, cursor uint64, pkg model.Package) (bool, error)
}

type repo struct {
	DB        *sqlx.DB
	batchSize uint
}

func (r *repo) DeletePackage(ctx context.Context, id uint64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) GetPackage(ctx context.Context, id uint64) (*model.Package, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) UpdatePackage(ctx context.Context, cursor uint64, pkg model.Package) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) Create(ctx context.Context, pkg model.Package) (*uint64, error) {

	query := sq.Insert("package").PlaceholderFormat(sq.Dollar).Columns("title", "created").Values(pkg.Title, pkg.CreatedAt).Suffix("RETURNING id")

	rows := query.QueryRowContext(ctx)

	id := uint64(0)
	err2 := rows.Scan(&id)
	if err2 != nil {
		return nil, err2
	}

	return &id, nil
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{DB: db, batchSize: batchSize}
}
