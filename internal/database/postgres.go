package database

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
)

/*
PostgreSQL Error Codes
https://www.postgresql.org/docs/16/errcodes-appendix.html
*/

// MustGetPgxPool get pgxpool or os.Exit(1)
func MustGetPgxPool(ctx context.Context) *pgxpool.Pool {

	log := slog.With("func", "database.MustGetPgxPool")

	dbpool, err := NewPgxPool(ctx)
	if err != nil {
		log.Warn("Failed init postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return dbpool
}

func NewPgxPool(ctx context.Context) (*pgxpool.Pool, error) {

	log := slog.With("func", "database.NewPgxPool")

	cfg := config.GetConfigInstance()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	// Эти параметры можно также задать в DSN
	/* дефолтные значения:
	   pool_max_conn_lifetime = time.Hour
	   pool_max_conn_idle_time = time.Minute * 30
	   pool_health_check_period = time.Minute
	   pool_max_conns = greater of 4 or runtime.NumCPU() если ядер больше 4	*/

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Warn("Error connecting to the database", slog.Any("error", err))
		return nil, fmt.Errorf("database.NewPgxPool: %w", err)
	}

	err = dbpool.Ping(ctx) // эта команда заменяет acquire + ping
	if err != nil {
		log.Warn("Could not ping database", slog.Any("error", err))
		return nil, fmt.Errorf("database.NewPgxPool: %w", err)
	}

	log.Info("successfully connected to database")
	return dbpool, nil
}
