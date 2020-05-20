package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/server/repository"
)

func TestBookRepoImpl_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: db}
	sql := regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)
	t.Run("sql error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectQuery(sql).WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("some-insert-error"))
		_, err = repo.Create(ctx, &repository.Book{Title: "some-title", Author: "some-author"})
		require.EqualError(t, err, "some-insert-error")
	})
	t.Run("sql success", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectQuery(sql).WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
		id, err := repo.Create(ctx, &repository.Book{Title: "some-title", Author: "some-author"})
		require.NoError(t, err)
		require.Equal(t, int64(999), id)
	})
}

func TestBookRepitory_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: db}
	t.Run("sql error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
			WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnError(fmt.Errorf("some-update-error"))

		require.EqualError(t,
			repo.Update(ctx, &repository.Book{ID: 888, Title: "new-title", Author: "new-author"}),
			"some-update-error",
		)
	})
	t.Run("sql success", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
			WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, repo.Update(ctx, &repository.Book{ID: 888, Title: "new-title", Author: "new-author"}))
	})
}

func TestBookRepoImpl_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: db}
	t.Run("sql error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
			WithArgs(666).
			WillReturnError(fmt.Errorf("some-delete-error"))
		err := repo.Delete(ctx, 666)
		require.EqualError(t, err, "some-delete-error")
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

func TestBookRepoImpl_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.BookRepoImpl{DB: db}
	t.Run("WHEN sql error", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
			WillReturnError(fmt.Errorf("some-list-error"))
		_, err := repo.Find(ctx)
		require.EqualError(t, err, "some-list-error")
	})
	t.Run("WHEN wrong dataset", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "tittle"}).
				AddRow(1, "one").
				AddRow(2, "two"))
		_, err := repo.Find(ctx)
		require.EqualError(t, err, "sql: expected 2 destination arguments in Scan, not 5")
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
