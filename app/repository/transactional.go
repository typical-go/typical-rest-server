package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"go.uber.org/dig"
)

const (
	// TxKey is key for tx
	TxKey key = iota
	// ErrKey is key for error of database
	ErrKey
)

// Transactional database
type Transactional struct {
	dig.In
	*sql.DB
}

type key int

// CommitMe to create begin transaction and return commit function to be deffered
func (t *Transactional) CommitMe(ctx *context.Context) func() {
	var (
		tx  *sql.Tx
		err error
	)
	if tx, err = t.DB.BeginTx(*ctx, nil); err != nil {
		*ctx = context.WithValue(*ctx, ErrKey, err)
		return func() {}
	}
	*ctx = context.WithValue(*ctx, TxKey, tx)
	return func() {
		if err = tx.Commit(); err != nil {
			*ctx = context.WithValue(*ctx, ErrKey, err)
		}
	}
}

// TxCtx return transaction from context if any or return t params
func TxCtx(ctx context.Context, t sq.BaseRunner) sq.BaseRunner {
	if tx, ok := ctx.Value(TxKey).(sq.BaseRunner); ok {
		return tx
	}
	return t
}

// ErrCtx return error from context
func ErrCtx(ctx context.Context) error {
	if err, ok := ctx.Value(ErrKey).(error); ok {
		return err
	}
	return nil
}
