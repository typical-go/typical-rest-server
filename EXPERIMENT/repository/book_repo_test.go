package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENT/repository"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/postgresdb"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
)

type bookRepoFn func(sqlmock.Sqlmock)

func createBookRepo(fn bookRepoFn) (repository.BookRepo, *sql.DB) {
	db, mock, _ := sqlmock.New()
	if fn != nil {
		fn(mock)
	}
	return repository.NewBookRepo(repository.BookRepoImpl{DB: db}), db
}

func TestBookRepoImpl_Create(t *testing.T) {
	testcases := []struct {
		TestName    string
		Book        *postgresdb.Book
		BookRepoFn  bookRepoFn
		Expected    int64
		ExpectedErr string
	}{
		{
			TestName:    "begin error",
			Book:        &postgresdb.Book{Title: "some-title", Author: "some-author"},
			ExpectedErr: "dbtxn: some-error",
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("some-error"))
			},
		},
		{
			TestName:    "insert error",
			Book:        &postgresdb.Book{Title: "some-title", Author: "some-author"},
			ExpectedErr: "some-error",
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some-error"))
			},
		},
		{
			Book: &postgresdb.Book{Title: "some-title", Author: "some-author"},
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
			},
			Expected: 999,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createBookRepo(tt.BookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			id, err := repo.Create(ctx, tt.Book)

			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, id)
			}

		})
	}
}

func TestBookRepoImpl_Update(t *testing.T) {
	testcases := []struct {
		TestName    string
		Book        *postgresdb.Book
		BookRepoFn  bookRepoFn
		Opt         dbkit.UpdateOption
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName:    "update error",
			Book:        &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:         dbkit.Equal(repository.BookTable.ID, 888),
			ExpectedErr: "dbtxn: begin-error",
			Expected:    -1,
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			TestName: "update error",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			ExpectedErr: "some-update-error",
			Expected:    -1,
		},
		{
			TestName: "bad update option",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt: dbkit.NewUpdateOption(func(b sq.UpdateBuilder) (sq.UpdateBuilder, error) {
				return b, errors.New("bad-option")
			}),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty author",
			Book:     &postgresdb.Book{Title: "new-title"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty title",
			Book:     &postgresdb.Book{Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createBookRepo(tt.BookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Update(ctx, tt.Book, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}

func TestBookRepoImpl_Patch(t *testing.T) {
	testcases := []struct {
		TestName    string
		Book        *postgresdb.Book
		BookRepoFn  bookRepoFn
		Opt         dbkit.UpdateOption
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName: "begin error",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
			ExpectedErr: "dbtxn: begin-error",
			Expected:    -1,
		},
		{
			TestName: "update error",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			ExpectedErr: "some-update-error",
			Expected:    -1,
		},

		{
			TestName: "bad update option",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt: dbkit.NewUpdateOption(func(b sq.UpdateBuilder) (sq.UpdateBuilder, error) {
				return b, errors.New("bad-option")
			}),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Book:     &postgresdb.Book{Title: "new-title", Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty author",
			Book:     &postgresdb.Book{Title: "new-title"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-title", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty title",
			Book:     &postgresdb.Book{Author: "new-author"},
			Opt:      dbkit.Equal(repository.BookTable.ID, 888),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET author = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createBookRepo(tt.BookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Patch(ctx, tt.Book, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}

func TestBookRepoImpl_Retrieve(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		TestName    string
		Opts        []dbkit.SelectOption
		Expected    []*postgresdb.Book
		ExpectedErr string
		BookRepoFn  bookRepoFn
	}{
		{
			TestName: "sql error",
			Opts:     []dbkit.SelectOption{},
			Expected: []*postgresdb.Book{},
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnError(errors.New("some-error"))
			},
			ExpectedErr: "some-error",
		},
		{
			TestName: "bad option",
			Opts: []dbkit.SelectOption{
				dbkit.NewSelectOption(func(b sq.SelectBuilder) (sq.SelectBuilder, error) {
					return b, errors.New("bad-option")
				}),
			},
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "author", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-author4", now, now).
						AddRow("1235", "some-title5", "some-author5", now, now),
					)
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Opts:     []dbkit.SelectOption{},
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "author", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-author4", now, now).
						AddRow("1235", "some-title5", "some-author5", now, now),
					)
			},
			Expected: []*postgresdb.Book{
				&postgresdb.Book{ID: 1234, Title: "some-title4", Author: "some-author4", UpdatedAt: now, CreatedAt: now},
				&postgresdb.Book{ID: 1235, Title: "some-title5", Author: "some-author5", UpdatedAt: now, CreatedAt: now},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createBookRepo(tt.BookRepoFn)
			defer db.Close()

			books, err := repo.Find(context.Background(), tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, books)
			}
		})
	}
}

func TestBookRepoImpl_Delete(t *testing.T) {
	testcases := []struct {
		TestName    string
		Opt         dbkit.DeleteOption
		BookRepoFn  bookRepoFn
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName:    "begin error",
			Opt:         dbkit.Equal("id", 666),
			ExpectedErr: "dbtxn: begin-error",
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			TestName:    "delete error",
			Opt:         dbkit.Equal("id", 666),
			ExpectedErr: "delete-error",
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(666).
					WillReturnError(errors.New("delete-error"))
			},
		},
		{
			TestName: "bad delete option",
			Opt: dbkit.NewDeleteOption(func(b sq.DeleteBuilder) (sq.DeleteBuilder, error) {
				return b, errors.New("bad-option")
			}),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Opt:      dbkit.Equal("id", 555),
			BookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createBookRepo(tt.BookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Delete(ctx, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}
