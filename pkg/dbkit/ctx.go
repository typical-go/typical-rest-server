package dbkit

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

const (
	// TxCtxKey is key for tx
	TxCtxKey key = iota
	// ErrCtxKey is key for error of database
	ErrCtxKey
)

type key int

// SetTxCtx to set tx into context
func SetTxCtx(ctx context.Context, tx sq.BaseRunner) context.Context {
	return context.WithValue(ctx, TxCtxKey, tx)
}

// SetErrCtx to set error into context
func SetErrCtx(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, ErrCtxKey, err)
}

// TxCtx return transaction from context if any or return t params
func TxCtx(ctx context.Context, t sq.BaseRunner) sq.BaseRunner {
	if tx, ok := ctx.Value(TxCtxKey).(sq.BaseRunner); ok {
		return tx
	}
	return t
}

// ErrCtx return error from context
func ErrCtx(ctx context.Context) error {
	if err, ok := ctx.Value(ErrCtxKey).(error); ok {
		return err
	}
	return nil
}
