package repository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/app/repository"
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
		require.EqualError(t, repository.ErrCtx(ctx),
			"all expectations were already fulfilled, call to database transaction Begin was not expected")
	})
	t.Run("WHEN commit error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		trx.CommitMe(&ctx)()
		require.EqualError(t, repository.ErrCtx(ctx),
			"all expectations were already fulfilled, call to Commit transaction was not expected")
	})
	t.Run("WHEN okay", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		mock.ExpectCommit()
		trx.CommitMe(&ctx)()
		require.NoError(t, repository.ErrCtx(ctx))
		require.NotNil(t, repository.TxCtx(ctx, nil))
	})
}
