package dbtxn_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
)

func TestRetrieve(t *testing.T) {
	testcases := []struct {
		testName        string
		ctx             context.Context
		expectedContext *dbtxn.Context
	}{
		{
			ctx:             nil,
			expectedContext: nil,
		},
		{
			ctx:             context.Background(),
			expectedContext: nil,
		},
		{
			ctx:             context.WithValue(context.Background(), dbtxn.ContextKey, "meh"),
			expectedContext: nil,
		},
		{
			ctx:             context.WithValue(context.Background(), dbtxn.ContextKey, &dbtxn.Context{}),
			expectedContext: &dbtxn.Context{},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expectedContext, dbtxn.Retrieve(tt.ctx))
		})
	}
}

func TestUse_When_Error(t *testing.T) {
	t.Run("Rollback when error", func(t *testing.T) {
		ctx := context.Background()
		defer dbtxn.Begin(&ctx)()

		db, mock, _ := sqlmock.New()
		mock.ExpectBegin()
		mock.ExpectRollback()

		txn, err := dbtxn.Use(ctx, db)
		require.NoError(t, err)
		require.True(t, txn.Txn())
		txn.SetError(errors.New("some-message"))
	})

	t.Run("Commit when no error", func(t *testing.T) {
		ctx := context.Background()
		defer dbtxn.Begin(&ctx)()

		db, mock, _ := sqlmock.New()
		mock.ExpectBegin()
		mock.ExpectCommit()

		txn, err := dbtxn.Use(ctx, db)
		require.NoError(t, err)
		require.True(t, txn.Txn())
	})
}
