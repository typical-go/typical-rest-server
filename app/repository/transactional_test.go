package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestTransactional(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	trx := repository.Transactional{
		DB: db,
	}
	t.Run("WHEN error occurred before commit", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		mock.ExpectRollback()
		commitFn := trx.CommitMe(&ctx)
		func(ctx context.Context) {
			dbkit.SetErrCtx(ctx, errors.New("unexpected-error"))
		}(ctx)
		commitFn()
		require.EqualError(t, dbkit.ErrCtx(ctx), "unexpected-error")
	})
	t.Run("WHEN panic occurred before commit", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		fn := trx.CommitMe(&ctx)
		func(ctx context.Context) { // service level
			defer fn()
			dbkit.SetErrCtx(ctx, fmt.Errorf("some-logic-error"))
			func(ctx context.Context) { // repository level
				panic("something-dangerous")
			}(ctx)
		}(ctx)
		require.EqualError(t, dbkit.ErrCtx(ctx), "something-dangerous")
	})
	t.Run("WHEN begin error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin().WillReturnError(errors.New("some-begin-error"))
		require.EqualError(t, trx.CommitMe(&ctx)(), "some-begin-error")
	})
	t.Run("WHEN commit error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.New("some-commit-error"))
		require.EqualError(t, trx.CommitMe(&ctx)(), "some-commit-error")
		require.NoError(t, dbkit.ErrCtx(ctx))
		require.NotNil(t, dbkit.TxCtx(ctx, nil))
	})
	t.Run("WHEN rolback error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(errors.New("some-rollback-error"))
		commitFn := trx.CommitMe(&ctx)
		func(ctx context.Context) {
			dbkit.SetErrCtx(ctx, errors.New("unexpected-error"))
		}(ctx)
		require.EqualError(t, commitFn(), "some-rollback-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "unexpected-error")
	})
}
