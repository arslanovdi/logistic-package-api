package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Repo is DAO for Template
type Repo interface {
	Describe(ctx context.Context, PackageID uint64) (*model.Package, error) // описание
	//Remove(ctx context.Context, packageID uint64) (bool, error)
	Create(ctx context.Context, pkg model.Package) (uint64, error)
	/*List(ctx context.Context, cursor uint64, limit uint64) ([]model.Package, error)
	Update(ctx context.Context, packageID uint64, pkg model.Package) error
	Get(ctx context.Context, cursor uint64) (model.Package, error)*/
}

type repo struct {
	DB        *sqlx.DB
	batchSize uint
}

func (r *repo) Create(ctx context.Context, pkg model.Package) (uint64, error) {
	query := sq.Insert("package").PlaceholderFormat(sq.Dollar).Columns("title", "created").Values(pkg.Title, pkg.CreatedAt).Suffix("RETURNING id")

	rows := query.QueryRowContext(ctx)

	id := uint64(0)
	err2 := rows.Scan(&id)
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{DB: db, batchSize: batchSize}
}

func (r *repo) Describe(ctx context.Context, templateID uint64) (*model.Package, error) {
	return nil, nil
}
