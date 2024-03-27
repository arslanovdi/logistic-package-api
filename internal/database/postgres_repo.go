package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/pkg/ctxutil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // Плэйсхолдер для Postgres

const stuckTimeout = 5 * time.Minute // время через которое залоченное событие считается зависшим и отправляется повторно

type repo struct {
	dbpool    *pgxpool.Pool
	batchSize uint
}

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

func (r *repo) Unlock(ctx context.Context, eventID []uint64) error {

	log := slog.With("func", "postgres.Unlock")

	query, args, err1 := psql.Update("package_events").
		Set("status", model.Unlocked).
		Where(sq.Eq{"id": eventID}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Unlock: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	tag, err2 := r.dbpool.Exec(ctx, query, args...)
	if err2 != nil {
		return fmt.Errorf("postgres.Unlock: %w", err2)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

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

func (r *repo) DeletePackage(ctx context.Context, id uint64) (bool, error) {

	log := slog.With("func", "postgres.DeletePackage")

	query, args, err1 := psql.Delete("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return false, fmt.Errorf("postgres.DeletePackage: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	queryEvent, argsEvent, err2 := psql.Insert("package_events").
		Columns("package_id", "type").
		Values(id, model.Removed).
		ToSql()

	if err2 != nil {
		return false, fmt.Errorf("postgres.DeletePackage: %w", err2)
	}

	log.Debug("queryEvent", slog.String("query", queryEvent), slog.Any("args", argsEvent))

	ctx = ctxutil.Detach(ctx)

	err3 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error {

		tag, err := r.dbpool.Exec(ctx, query, args...)
		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 { // Получаем количество обновленных строк
			return model.ErrNotFound
		}

		_, err = r.dbpool.Exec(ctx, queryEvent, argsEvent...)
		if err != nil {
			return err
		}

		return nil
	})

	if err3 != nil {
		return false, fmt.Errorf("postgres.DeletePackage: %w", err3)
	}

	log.Debug("Package deleted", slog.Any("id", id))

	return true, nil
}

func (r *repo) GetPackage(ctx context.Context, id uint64) (*model.Package, error) {

	log := slog.With("func", "postgres.GetPackage")

	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.GetPackage: %w", err1)
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
		return nil, fmt.Errorf("postgres.GetPackage: %w", err2)
	}

	log.Debug("GetPackage", slog.Any("pkg", pkg))

	return &pkg, nil
}

func (r *repo) ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {

	log := slog.With("func", "postgres.ListPackages")

	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.GtOrEq{"id": offset}).
		Where(sq.Lt{"id": offset + limit}).
		OrderBy("id ASC").
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.ListPackages: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	rows, _ := r.dbpool.Query(ctx, query, args...) // Ошибка игнорируется, так как она обрабатывается в CollectRows
	defer rows.Close()

	var packages []model.Package
	packages, err2 := pgx.CollectRows(rows, pgx.RowToStructByName[model.Package]) // десериализуем в слайс структуру
	if err2 != nil {
		if errors.Is(err2, pgx.ErrNoRows) {
			log.Debug("no rows found")
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("postgres.ListPackages: %w", err2)
	}

	log.Debug("packages listed", slog.Uint64("offset", offset), slog.Uint64("limit", limit))

	return packages, nil
}

// UpdatePackage - update package by id in Postgres
func (r *repo) UpdatePackage(ctx context.Context, pkg model.Package) (bool, error) {

	log := slog.With("func", "postgres.UpdatePackage")

	query, args, err1 := psql.Update("package").
		Set("weight", pkg.Weight).
		Set("title", pkg.Title).
		Set("updated", pkg.Updated).
		Where(sq.Eq{"id": pkg.ID}).
		Suffix("RETURNING created, removed").
		ToSql()
	if err1 != nil {
		return false, fmt.Errorf("postgres.UpdatePackage: %w", err1)
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

		_, err = r.dbpool.Exec(ctx, queryEvent, argsEvent...)
		if err != nil {
			return err
		}

		return nil
	})

	if err2 != nil {
		return false, fmt.Errorf("postgres.UpdatePackage: %w", err2)
	}

	log.Debug("package updated", slog.String("package", pkg.String()))

	return true, nil
}

// Create - create new package in Postgres
func (r *repo) Create(ctx context.Context, pkg model.Package) (*uint64, error) {

	log := slog.With("func", "postgres.Create")

	query, args, err1 := psql.Insert("package").
		Columns("weight", "title", "created").
		Values(pkg.Weight, pkg.Title, pkg.Created).
		Suffix("RETURNING id").
		ToSql()
	if err1 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	err2 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error { // Запускаем транзакцию

		err := r.dbpool.QueryRow(ctx, query, args...).Scan(&pkg.ID)
		if err != nil {
			return err
		}

		pkgJSON, err := json.Marshal(pkg)
		if err != nil {
			return err
		}

		queryEvent, argsEvent, err := psql.Insert("package_events").
			Columns("package_id", "type", "payload").
			Values(pkg.ID, model.Created, pkgJSON).
			ToSql()
		if err != nil {
			return err
		}

		log.Debug("queryEvent", slog.String("query", queryEvent), slog.Any("args", argsEvent))

		_, err = r.dbpool.Exec(ctx, queryEvent, argsEvent...)
		if err != nil {
			return err
		}

		return nil
	})

	if err2 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err2)
	}

	log.Debug("package created", slog.String("package", pkg.String()))

	return &pkg.ID, nil
}

// NewPostgresRepo returns Postgres implementation of service.Repo and retranslator.EventRepo
func NewPostgresRepo(dbpool *pgxpool.Pool, batchSize uint) *repo {
	return &repo{
		dbpool:    dbpool,
		batchSize: batchSize,
	}
}
