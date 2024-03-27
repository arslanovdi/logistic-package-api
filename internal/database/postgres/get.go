package postgres

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/pkg/ctxutil"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

func (r *repo) Get(ctx context.Context, id uint64) (*model.Package, error) {

	log := slog.With("func", "postgres.Get")

	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.Get: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	rows, _ := r.dbpool.Query(ctx, query, args...)
	defer rows.Close()

	pkg, err2 := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Package])
	if err2 != nil {
		if errors.Is(err2, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("postgres.Get: %w", err2)
	}

	log.Debug("Get", slog.Any("pkg", pkg))

	return &pkg, nil
}
