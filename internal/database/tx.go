package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WithTxFunc func(ctx context.Context, tx *sqlx.Tx) error

// WithTx обертка, автоматически попытается зароллбэчить транзакцию при ошибке
func WithTx(ctx context.Context, db *sqlx.DB, fn WithTxFunc) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{ // запуск транзакции с опциями
		Isolation: sql.LevelRepeatableRead, // можно пробрасывать уровень изоляции при помощи паттерна функциональных опций
		ReadOnly:  false,
	})
	if err != nil {
		return errors.Wrap(err, "db.BeginTx()")
	}
	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, "Tx.Rollback()")
		}
		return errors.Wrap(err, "Tx.WithTxFunc()")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "Tx.Commit()")
	}
	return nil
}
