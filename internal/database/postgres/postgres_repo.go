// Package postgres - Postgres implementation of service.Repo and repo.EventRepo
package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // Плэйсхолдер для Postgres

const stuckTimeout = 5 * time.Minute // время через которое залоченное событие считается зависшим и отправляется повторно

// Repo - Postgres implementation of service.Repo and repo.EventRepo
type Repo struct {
	dbpool *pgxpool.Pool
	//batchSize uint
}

// NewPostgresRepo returns Postgres implementation of service.Repo and repo.EventRepo
func NewPostgresRepo(dbpool *pgxpool.Pool) *Repo {
	return &Repo{
		dbpool: dbpool,
		//batchSize: batchSize,
	}
}
