package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/app/repository"
)

func TestBookRepo_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	repo := repository.NewBookRepo(repository.CachedBookRepoImpl{
		BookRepoImpl: repository.BookRepoImpl{DB: db},
		Redis:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
	})
	ctx0 := context.Background()
	insertSQL := regexp.QuoteMeta(`INSERT INTO books (title,author) VALUES ($1,$2) RETURNING "id"`)
	t.Run("sql error", func(t *testing.T) {
		mock.ExpectQuery(insertSQL).WithArgs("some-title", "some-author").
			WillReturnError(fmt.Errorf("some-insert-error"))
		_, err = repo.Insert(ctx0, repository.Book{Title: "some-title", Author: "some-author"})
		require.EqualError(t, err, "some-insert-error")
	})
	t.Run("sql success", func(t *testing.T) {
		mock.ExpectQuery(insertSQL).WithArgs("some-title", "some-author").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
		id, err := repo.Insert(ctx0, repository.Book{Title: "some-title", Author: "some-author"})
		require.NoError(t, err)
		require.Equal(t, int64(999), id)
	})
}

func TestBookRepitory_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	repo := repository.NewBookRepo(repository.CachedBookRepoImpl{
		BookRepoImpl: repository.BookRepoImpl{DB: db},
		Redis:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
	})
	ctx0 := context.Background()
	updateSQL := regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)
	t.Run("sql error", func(t *testing.T) {
		mock.ExpectExec(updateSQL).WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnError(fmt.Errorf("some-update-error"))
		err = repo.Update(ctx0, repository.Book{ID: 888, Title: "new-title", Author: "new-author"})
		require.EqualError(t, err, "some-update-error")
	})
	t.Run("sql success", func(t *testing.T) {
		mock.ExpectExec(updateSQL).WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err = repo.Update(ctx0, repository.Book{ID: 888, Title: "new-title", Author: "new-author"})
		require.NoError(t, err)
	})
}

func TestBookRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	repo := repository.NewBookRepo(repository.CachedBookRepoImpl{
		BookRepoImpl: repository.BookRepoImpl{DB: db},
		Redis:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
	})
	ctx0 := context.Background()
	deleteSQL := regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)
	t.Run("sql error", func(t *testing.T) {
		mock.ExpectExec(deleteSQL).WithArgs(666).
			WillReturnError(fmt.Errorf("some-delete-error"))
		err := repo.Delete(ctx0, 666)
		require.EqualError(t, err, "some-delete-error")
	})
	t.Run("sql success", func(t *testing.T) {
		mock.ExpectExec(deleteSQL).WithArgs(555).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.Delete(ctx0, 555)
		require.NoError(t, err)
	})
}

func TestBookRepitory_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	repo := repository.NewBookRepo(repository.CachedBookRepoImpl{
		BookRepoImpl: repository.BookRepoImpl{DB: db},
		Redis:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
	})
	ctx0 := context.Background()
	querySQL := regexp.QuoteMeta(`SELECT id, title, author, updated_at, created_at FROM books WHERE id = $1`)
	t.Run("WHEN sql error", func(t *testing.T) {
		mock.ExpectQuery(querySQL).WithArgs(123).
			WillReturnError(fmt.Errorf("some-find-error"))
		_, err := repo.Find(ctx0, 123)
		require.EqualError(t, err, "some-find-error")
	})
	t.Run("WHEN cache available", func(t *testing.T) {

	})
	t.Run("WHEN okay", func(t *testing.T) {
		expected := &repository.Book{
			ID:        123,
			Title:     "some-title",
			Author:    "some-author",
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}
		mock.ExpectQuery(querySQL).WithArgs(123).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "updated_at", "created_at"}).
				AddRow(expected.ID, expected.Title, expected.Author, expected.UpdatedAt, expected.CreatedAt))
		book, err := repo.Find(ctx0, 123)
		require.NoError(t, err)
		require.Equal(t, expected, book)
		data, _ := testRedis.Get("BOOK:FIND:123")
		defer testRedis.FlushAll()
		require.Equal(t, `{"id":123,"title":"some-title","author":"some-author"}`, data)
	})
}

func TestBookRepo_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	repo := repository.NewBookRepo(repository.CachedBookRepoImpl{
		BookRepoImpl: repository.BookRepoImpl{DB: db},
		Redis:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
	})
	ctx0 := context.Background()
	listSQL := `SELECT id, title, author, updated_at, created_at FROM books`
	t.Run("WHEN sql error", func(t *testing.T) {
		mock.ExpectQuery(listSQL).WillReturnError(fmt.Errorf("some-list-error"))
		_, err := repo.List(ctx0)
		require.EqualError(t, err, "some-list-error")
	})
	t.Run("WHEN wrong dataset", func(t *testing.T) {
		mock.ExpectQuery(listSQL).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
			AddRow(1, "one").
			AddRow(2, "two"))
		_, err := repo.List(ctx0)
		require.EqualError(t, err, "sql: expected 2 destination arguments in Scan, not 5")
	})
	t.Run("WHEN okay", func(t *testing.T) {
		expecteds := []*repository.Book{
			&repository.Book{ID: 1234, Title: "some-title4", Author: "some-author4", UpdatedAt: time.Now(), CreatedAt: time.Now()},
			&repository.Book{ID: 1235, Title: "some-title5", Author: "some-author5", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		}
		rows := sqlmock.NewRows([]string{"id", "title", "author", "updated_at", "created_at"})
		for _, expected := range expecteds {
			rows.AddRow(expected.ID, expected.Title, expected.Author, expected.UpdatedAt, expected.CreatedAt)
		}
		mock.ExpectQuery(listSQL).WillReturnRows(rows)
		books, err := repo.List(ctx0)
		require.NoError(t, err)
		require.Equal(t, expecteds, books)
		data, _ := testRedis.Get("BOOK:LIST")
		defer testRedis.FlushAll()
		require.Equal(t, `[{"id":1234,"title":"some-title4","author":"some-author4"},{"id":1235,"title":"some-title5","author":"some-author5"}]`, data)
	})
	t.Run("WHEN cache available", func(t *testing.T) {
		testRedis.Set("BOOK:LIST", `[{"id":1234,"title":"some-title4","author":"some-author4"},{"id":1235,"title":"some-title5","author":"some-author5"}]`)
		// defer testRedis.FlushAll()
		books, err := repo.List(ctx0)
		require.NoError(t, err)
		require.Equal(t, []*repository.Book{
			&repository.Book{ID: 1234, Title: "some-title4", Author: "some-author4"},
			&repository.Book{ID: 1235, Title: "some-title5", Author: "some-author5"},
		}, books)

	})

}
