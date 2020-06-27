package dbtxn_test

import (
	"context"
	"database/sql"
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

func TestUse(t *testing.T) {
	testcases := []struct {
		testName    string
		ctx         context.Context
		db          *sql.DB
		expected    *dbtxn.Handler
		expectedErr string
	}{
		{
			ctx:         nil,
			expectedErr: "dbtxn: missing context.Context",
		},
		{
			testName: "non transactional",
			db:       &sql.DB{},
			ctx:      context.Background(),
			expected: &dbtxn.Handler{DB: &sql.DB{}},
		},
		{
			testName: "begin error",
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
				return db
			}(),
			ctx:         context.WithValue(context.Background(), dbtxn.ContextKey, &dbtxn.Context{}),
			expectedErr: "dbtxn: begin-error",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			handler, err := dbtxn.Use(tt.ctx, tt.db)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, handler)
			}
		})
	}
}

func TestUse_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectBegin()

	ctx := context.WithValue(context.Background(), dbtxn.ContextKey, &dbtxn.Context{})
	handler, err := dbtxn.Use(ctx, db)

	require.NoError(t, err)
	require.Equal(t, handler.DB, handler.Context.Tx)
}

func TestContext_Commit(t *testing.T) {

	t.Run("expect rollback when error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectBegin()
		mock.ExpectRollback()

		tx, _ := db.Begin()
		c := &dbtxn.Context{Tx: tx}
		c.Err = errors.New("some-error")

		require.NoError(t, c.Commit())
	})

	t.Run("expect commit when no error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectBegin()
		mock.ExpectCommit()

		tx, _ := db.Begin()
		c := &dbtxn.Context{Tx: tx}

		require.NoError(t, c.Commit())
	})

}
