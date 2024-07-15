package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/pkg/ctxutil"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// Update - update package by id in database
func (r *Repo) Update(ctx context.Context, pkg model.Package) error {

	log := slog.With("func", "postgres.Update")

	query, args, err1 := psql.Update("package").
		Set("weight", pkg.Weight).
		Set("title", pkg.Title).
		Set("updated", pkg.Updated).
		Where(sq.Eq{"id": pkg.ID}).
		Suffix("RETURNING created, removed").
		ToSql()
	if err1 != nil {
		return fmt.Errorf("postgres.Update: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	err2 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error {

		err := tx.QueryRow(ctx, query, args...).Scan(&pkg.Created, &pkg.Removed)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrNotFound
			}
			return err
		}

		pkgJSON, err := json.Marshal(pkg)
		if err != nil {
			return err
		}

		queryEvent, argsEvent, err := psql.Insert("package_events").
			Columns("package_id", "type", "payload").
			Values(pkg.ID, model.Updated, pkgJSON).
			ToSql()
		if err != nil {
			return err
		}

		log.Debug("queryEvent", slog.String("query", queryEvent), slog.Any("args", argsEvent))

		_, err = tx.Exec(ctx, queryEvent, argsEvent...)
		//_, err = r.dbpool.Exec(ctx, queryEvent, argsEvent...)
		if err != nil {
			return err
		}

		return nil
	})

	if err2 != nil {
		return fmt.Errorf("postgres.Update: %w", err2)
	}

	log.Debug("package updated", slog.String("package", pkg.String()))

	return nil
}
