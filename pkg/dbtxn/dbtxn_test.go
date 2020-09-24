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
		TestName        string
		Ctx             context.Context
		ExpectedContext *dbtxn.Context
	}{
		{
			Ctx:             nil,
			ExpectedContext: nil,
		},
		{
			Ctx:             context.Background(),
			ExpectedContext: nil,
		},
		{
			Ctx:             context.WithValue(context.Background(), dbtxn.ContextKey, "meh"),
			ExpectedContext: nil,
		},
		{
			Ctx:             context.WithValue(context.Background(), dbtxn.ContextKey, &dbtxn.Context{}),
			ExpectedContext: &dbtxn.Context{},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.ExpectedContext, dbtxn.Find(tt.Ctx))
		})
	}
}

func TestUse(t *testing.T) {
	testcases := []struct {
		TestName    string
		Ctx         context.Context
		DB          *sql.DB
		Expected    *dbtxn.Handler
		ExpectedErr string
	}{
		{
			Ctx:         nil,
			ExpectedErr: "dbtxn: missing context.Context",
		},
		{
			TestName: "non transactional",
			DB:       &sql.DB{},
			Ctx:      context.Background(),
			Expected: &dbtxn.Handler{DB: &sql.DB{}},
		},
		{
			TestName: "begin error",
			DB: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
				return db
			}(),
			Ctx:         context.WithValue(context.Background(), dbtxn.ContextKey, &dbtxn.Context{}),
			ExpectedErr: "dbtxn: begin-error",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			handler, err := dbtxn.Use(tt.Ctx, tt.DB)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, handler)
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

func TestSetError(t *testing.T) {
	ctx := context.Background()
	dbtxn.Begin(&ctx)

	db, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	handler, err := dbtxn.Use(ctx, db)
	require.NoError(t, err)

	handler.SetError(errors.New("some-error"))
	require.EqualError(t, dbtxn.Error(ctx), "some-error")
}
