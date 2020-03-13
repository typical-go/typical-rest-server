package dbkit

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"go.uber.org/dig"
)

// Transactional database
type Transactional struct {
	dig.In
	*typpostgres.DB
}

// CommitMe to create begin transaction and return commit function to be deffered
func (t *Transactional) CommitMe(ctx *context.Context) func() error {
	var (
		tx  *sql.Tx
		err error
	)
	*ctx = CtxWithTxo(*ctx)
	if tx, err = t.DB.BeginTx(*ctx, nil); err != nil {
		return func() error {
			if r := recover(); r != nil {
				SetErrCtx(*ctx, fmt.Errorf("%v", r))
			}
			return err
		}
	}
	SetTxCtx(*ctx, tx)
	return func() error {
		if r := recover(); r != nil {
			SetErrCtx(*ctx, fmt.Errorf("%v", r))
			return tx.Rollback()
		}
		if err := ErrCtx(*ctx); err != nil {
			return tx.Rollback()
		}
		return tx.Commit()
	}
}

// CancelMe is store error to context to trigger the rollback mechanism
func (t *Transactional) CancelMe(ctx context.Context, err error) error {
	return SetErrCtx(ctx, err)
}
