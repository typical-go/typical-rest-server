package repository_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/app/repository"
)

func TestBookRepoImpl_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: typpostgres.NewDB(db)}
	sql := regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)
	t.Run("sql error", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectQuery(sql).WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("some-insert-error"))
		_, err = repo.Create(ctx, &repository.Book{Title: "some-title", Author: "some-author"})
		require.EqualError(t, err, "some-insert-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "some-insert-error")
	})
	t.Run("sql success", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectQuery(sql).WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
		book, err := repo.Create(ctx, &repository.Book{Title: "some-title", Author: "some-author"})
		require.NoError(t, err)
		require.Equal(t, int64(999), book.ID)
	})
}

func TestBookRepitory_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: typpostgres.NewDB(db)}
	t.Run("sql error", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		expectFindOneBook(mock, 888, &repository.Book{ID: 888})
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
			WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnError(fmt.Errorf("some-update-error"))
		_, err = repo.Update(ctx, &repository.Book{ID: 888, Title: "new-title", Author: "new-author"})
		require.EqualError(t, err, "some-update-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "some-update-error")
	})
	t.Run("sql success", func(t *testing.T) {
		ctx := context.Background()
		expectFindOneBook(mock, 888, &repository.Book{ID: 888})
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
			WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = repo.Update(ctx, &repository.Book{ID: 888, Title: "new-title", Author: "new-author"})
		require.NoError(t, err)
	})
}

func TestBookRepoImpl_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: typpostgres.NewDB(db)}
	t.Run("sql error", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
			WithArgs(666).
			WillReturnError(fmt.Errorf("some-delete-error"))
		err := repo.Delete(ctx, 666)
		require.EqualError(t, err, "some-delete-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "some-delete-error")
	})
	t.Run("sql success", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
			WithArgs(555).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.Delete(ctx, 555)
		require.NoError(t, err)
	})
}

func TestBookRepitory_FindOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: typpostgres.NewDB(db)}
	t.Run("WHEN sql error", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, author, updated_at, created_at FROM books WHERE id = $1`)).
			WithArgs(123).
			WillReturnError(errors.New("some-find-error"))
		_, err := repo.FindOne(ctx, 123)
		require.EqualError(t, err, "some-find-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "some-find-error")
	})
	t.Run("WHEN result set", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, author, updated_at, created_at FROM books WHERE id = $1`)).
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).
				AddRow("some-id", "some-title", "some-author"))
		_, err := repo.FindOne(ctx, 123)
		require.EqualError(t, err, "sql: expected 3 destination arguments in Scan, not 5")
		require.EqualError(t, dbkit.ErrCtx(ctx), "sql: expected 3 destination arguments in Scan, not 5")
	})
	t.Run("WHEN okay", func(t *testing.T) {
		ctx := context.Background()
		expected := &repository.Book{
			ID:        123,
			Title:     "some-title",
			Author:    "some-author",
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}
		expectFindOneBook(mock, 123, expected)
		book, err := repo.FindOne(ctx, 123)
		require.NoError(t, err)
		require.Equal(t, expected, book)
	})
}

func TestBookRepoImpl_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: typpostgres.NewDB(db)}
	t.Run("WHEN sql error", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
			WillReturnError(fmt.Errorf("some-list-error"))
		_, err := repo.Find(ctx)
		require.EqualError(t, err, "some-list-error")
		require.EqualError(t, dbkit.ErrCtx(ctx), "some-list-error")
	})
	t.Run("WHEN wrong dataset", func(t *testing.T) {
		ctx := dbkit.CtxWithTxo(context.Background())
		mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "tittle"}).
				AddRow(1, "one").
				AddRow(2, "two"))
		_, err := repo.Find(ctx)
		require.EqualError(t, err, "sql: expected 2 destination arguments in Scan, not 5")
		require.EqualError(t, dbkit.ErrCtx(ctx), "sql: expected 2 destination arguments in Scan, not 5")
	})
	t.Run("WHEN okay", func(t *testing.T) {
		ctx := context.Background()
		expecteds := []*repository.Book{
			&repository.Book{ID: 1234, Title: "some-title4", Author: "some-author4", UpdatedAt: time.Now(), CreatedAt: time.Now()},
			&repository.Book{ID: 1235, Title: "some-title5", Author: "some-author5", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		}
		rows := sqlmock.NewRows([]string{"id", "title", "author", "updated_at", "created_at"})
		for _, expected := range expecteds {
			rows.AddRow(expected.ID, expected.Title, expected.Author, expected.UpdatedAt, expected.CreatedAt)
		}
		mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
			WillReturnRows(rows)
		books, err := repo.Find(ctx)
		require.NoError(t, err)
		require.Equal(t, expecteds, books)
	})
}

func expectFindOneBook(mock sqlmock.Sqlmock, id driver.Value, result *repository.Book) *sqlmock.ExpectedQuery {
	return mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, author, updated_at, created_at FROM books WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "updated_at", "created_at"}).
			AddRow(result.ID, result.Title, result.Author, result.UpdatedAt, result.CreatedAt))
}
