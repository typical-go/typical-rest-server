package dbkit

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const (
	// TxoKey is key for Txo
	TxoKey key = iota
)

type key int

// Txo stand of transaction object
type Txo struct {
	tx  sq.BaseRunner
	err error
}

// CtxWithTxo return context with txo
func CtxWithTxo(parent context.Context) context.Context {
	return context.WithValue(parent, TxoKey, &Txo{})
}

func getTxo(ctx context.Context) *Txo {
	if err, ok := ctx.Value(TxoKey).(*Txo); ok {
		return err
	}
	return nil
}

// SetTxCtx to set tx in ctx
func SetTxCtx(ctx context.Context, tx sq.BaseRunner) error {
	if txo := getTxo(ctx); txo != nil {
		txo.tx = tx
		return nil
	}
	return errors.New("Context have no TXO")
}

// SetErrCtx to set tx in ctx
func SetErrCtx(ctx context.Context, err error) error {
	if txo := getTxo(ctx); txo != nil {
		txo.err = err
		return nil
	}
	return errors.New("Context have no TXO")
}

// TxCtx return transaction from context if any or return t params
func TxCtx(ctx context.Context, t sq.BaseRunner) sq.BaseRunner {
	if txo := getTxo(ctx); txo != nil {
		if txo.tx != nil {
			return txo.tx
		}
	}
	return t
}

// ErrCtx return error from context
func ErrCtx(ctx context.Context) error {
	if txo := getTxo(ctx); txo != nil {
		if txo.err != nil {
			return txo.err
		}
	}
	return nil
}
