package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"go.uber.org/dig"
)

// Transactional database
type Transactional struct {
	dig.In
	*sql.DB
}

// CommitMe to create begin transaction and return commit function to be deffered
func (t *Transactional) CommitMe(ctx *context.Context) func() error {
	var (
		tx  *sql.Tx
		err error
	)
	*ctx = dbkit.CtxWithTxo(*ctx)
	if tx, err = t.DB.BeginTx(*ctx, nil); err != nil {
		return func() error {
			if r := recover(); r != nil {
				dbkit.SetErrCtx(*ctx, fmt.Errorf("%v", r))
			}
			return err
		}
	}
	dbkit.SetTxCtx(*ctx, tx)
	return func() error {
		if r := recover(); r != nil {
			dbkit.SetErrCtx(*ctx, fmt.Errorf("%v", r))
			return tx.Rollback()
		}
		if err := dbkit.ErrCtx(*ctx); err != nil {
			return tx.Rollback()
		}
		return tx.Commit()
	}
}

// CancelMe is store error to context to trigger the rollback mechanism
func (t *Transactional) CancelMe(ctx context.Context, err error) error {
	return dbkit.SetErrCtx(ctx, err)
}
