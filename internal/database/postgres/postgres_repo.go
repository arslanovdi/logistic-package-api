package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // Плэйсхолдер для Postgres

const stuckTimeout = 5 * time.Minute // время через которое залоченное событие считается зависшим и отправляется повторно

type repo struct {
	dbpool    *pgxpool.Pool
	batchSize uint
}

// NewPostgresRepo returns Postgres implementation of service.Repo and retranslator.EventRepo
func NewPostgresRepo(dbpool *pgxpool.Pool, batchSize uint) *repo {
	return &repo{
		dbpool:    dbpool,
		batchSize: batchSize,
	}
}
