package database

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package-api/pkg/ctxutil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // Плэйсхолдер для Postgres

type repo struct {
	dbpool    *pgxpool.Pool
	batchSize uint
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

	ctx = ctxutil.Detach(ctx)

	tag, err2 := r.dbpool.Exec(ctx, query, args...)
	if err2 != nil {
		return false, fmt.Errorf("postgres.DeletePackage: %w", err2)
	}

	if tag.RowsAffected() == 0 { // Получаем количество обновленных строк
		return false, model.ErrNotFound
	}

	log.Debug("Package deleted", slog.Any("id", id))

	return true, nil
}

func (r *repo) GetPackage(ctx context.Context, id uint64) (*model.Package, error) {

	log := slog.With("func", "postgres.GetPackage")

	query, args, err1 := psql.Select("id", "weight", "title", "createdAt").
		From("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.GetPackage: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	rows, _ := r.dbpool.Query(ctx, query, args...)

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

	query, args, err1 := psql.Select("id", "weight", "title", "createdAt").
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
func (r *repo) UpdatePackage(ctx context.Context, cursor uint64, pkg model.Package) (bool, error) {

	log := slog.With("func", "postgres.UpdatePackage")

	query, args, err1 := psql.Update("package").
		Set("weight", pkg.Weight).
		Set("title", pkg.Title).
		Set("createdAt", pkg.CreatedAt).
		Where(sq.Eq{"id": cursor}).
		ToSql()

	if err1 != nil {
		return false, fmt.Errorf("postgres.UpdatePackage: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	result, err2 := r.dbpool.Exec(ctx, query, args...)
	if err2 != nil {
		return false, fmt.Errorf("postgres.UpdatePackage: %w", err2)
	}

	if result.RowsAffected() == 0 { // Получаем количество обновленных строк
		return false, model.ErrNotFound
	}

	log.Debug("package updated", slog.Uint64("id", cursor), slog.String("package", pkg.String()))

	return true, nil
}

// Create - create new package in Postgres
func (r *repo) Create(ctx context.Context, pkg model.Package) (id *uint64, err error) {

	log := slog.With("func", "postgres.Create")

	query, args, err1 := psql.Insert("package").
		Columns("weight", "title", "createdAt").
		Values(pkg.Weight, pkg.Title, pkg.CreatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	err2 := r.dbpool.QueryRow(ctx, query, args...).Scan(&id)

	if err2 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err2)
	}

	log.Debug("package created", slog.Uint64("id", *id), slog.String("package", pkg.String()))

	return id, nil
}

// NewPostgresRepo returns Postgres implementation of service.Repo
func NewPostgresRepo(dbpool *pgxpool.Pool, batchSize uint) service.Repo {
	return &repo{
		dbpool:    dbpool,
		batchSize: batchSize,
	}
}
