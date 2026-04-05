package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
)

type contextKey string

const TxKey contextKey = "tx"

type bunTransactionManager struct {
	db *bun.DB
}

func (m *bunTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, TxKey, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit()
}

func GetDB(ctx context.Context, defaultDB *bun.DB) bun.IDB {
	if tx, ok := ctx.Value(TxKey).(bun.Tx); ok {
		return tx
	}

	return defaultDB
}
