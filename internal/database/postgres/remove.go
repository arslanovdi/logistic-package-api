package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"log/slog"
)

func (r *repo) Remove(ctx context.Context, eventIDs []uint64) error {

	log := slog.With("func", "postgres.Remove")

	query, args, err1 := psql.Delete("package_events").
		Where(sq.Eq{"id": eventIDs}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Remove: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	tag, err2 := r.dbpool.Exec(ctx, query, args...)
	if err2 != nil {
		return fmt.Errorf("postgres.Remove: %w", err2)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
