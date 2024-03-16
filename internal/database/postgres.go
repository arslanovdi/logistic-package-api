package database

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// NewPostgres returns DB
func NewPostgres(dsn, driver string) (*sqlx.DB, error) {

	log := slog.With("func", "database.NewPostgres")

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Error("failed to create database connection", slog.Any("error", err))

		return nil, err
	}

	// need to uncomment for homework-4
	// if err = db.Ping(); err != nil {
	// 	log.Error().Err(err).Msgf("failed ping the database")

	// 	return nil, err
	// }

	return db, nil
}
