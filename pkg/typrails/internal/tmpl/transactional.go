package tmpl

// Transactional template
const Transactional = `package repository

import (
	"context"
	"database/sql"

	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"go.uber.org/dig"
)

// Transactional database
type Transactional struct {
	dig.In
	*sql.DB
}

// CommitMe to create begin transaction and return commit function to be deffered
func (t *Transactional) CommitMe(ctx *context.Context) func() {
	var (
		tx  *sql.Tx
		err error
	)
	if tx, err = t.DB.BeginTx(*ctx, nil); err != nil {
		*ctx = typrails.SetErrCtx(*ctx, err)
		return func() {}
	}
	*ctx = typrails.SetTxCtx(*ctx, tx)
	return func() {
		if err = tx.Commit(); err != nil {
			*ctx = typrails.SetErrCtx(*ctx, err)
		}
	}
}
`

// TransactionalTest template
const TransactionalTest = `package repository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/restserver/repository"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
)

func TestTransactional(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	trx := repository.Transactional{
		DB: db,
	}
	t.Run("WHEN begin error", func(t *testing.T) {
		ctx := context.Background()
		trx.CommitMe(&ctx)
		require.EqualError(t, dbkit.ErrCtx(ctx),
			"all expectations were already fulfilled, call to database transaction Begin was not expected")
	})
	t.Run("WHEN commit error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		trx.CommitMe(&ctx)()
		require.EqualError(t, dbkit.ErrCtx(ctx),
			"all expectations were already fulfilled, call to Commit transaction was not expected")
	})
	t.Run("WHEN okay", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		mock.ExpectCommit()
		trx.CommitMe(&ctx)()
		require.NoError(t, dbkit.ErrCtx(ctx))
		require.NotNil(t, dbkit.TxCtx(ctx, nil))
	})
}
`
