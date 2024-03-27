package postgres

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

func (r *repo) Lock(ctx context.Context, n uint64) ([]model.PackageEvent, error) {

	log := slog.With("func", "postgres.Lock")

	query, args, err1 := psql.Update("package_events").
		Set("status", model.Locked).
		Set("updated", time.Now()).
		Where(sq.Or{
			sq.NotEq{"status": model.Locked}, // если статус не залочен
			sq.And{
				//sq.NotEq{"updated": nil},	// TODO протестировать
				sq.LtOrEq{"updated": time.Now().Add(-stuckTimeout)}, // если событие зависло

			},
		}).
		OrderBy("id").
		Limit(n).
		Suffix("RETURNING *").
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.Lock: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	var events []model.PackageEvent

	err2 := pgx.BeginTxFunc(ctx, r.dbpool, pgx.TxOptions{IsoLevel: "serializable"}, func(tx pgx.Tx) error {

		rows, _ := r.dbpool.Query(ctx, query, args...)
		defer rows.Close()

		var err error
		events, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.PackageEvent]) // десериализуем в слайс структуру
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Debug("no rows found")
				return model.ErrNotFound
			}
			return err
		}

		return nil
	})

	if err2 != nil {
		return nil, fmt.Errorf("postgres.Lock: %w", err2)
	}

	return events, nil
}
