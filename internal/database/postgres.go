package database

import (
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/config"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// NewPostgres returns DB
func NewPostgres() (*sqlx.DB, error) {

	log := slog.With("func", "database.NewPostgres")

	cfg := config.GetConfigInstance()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := sqlx.Open(cfg.Database.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("database.NewPostgres.Open: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database.NewPostgres.Ping: %w", err)
	}

	log.Info("successfully connected to database")
	return db, nil
}
