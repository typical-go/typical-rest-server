package repository

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/imantung/go-helper/timekit"
	"github.com/stretchr/testify/require"
)

func TestBookRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows(bookColumns).
				AddRow(1, "title1", "author1", timekit.UTC("2019-02-20T10:00:00-05:00")),
		)
	mock.ExpectQuery("SELECT").WithArgs(9999).WillReturnError(fmt.Errorf("some-error"))

	bookRepository := NewBookRepository(db)

	t.Run("return rows", func(t *testing.T) {
		book, err := bookRepository.Get(1)
		require.NoError(t, err)
		require.Equal(t, book, Book{1, "title1", "author1", timekit.UTC("2019-02-20T10:00:00-05:00")})
	})

	t.Run("return error", func(t *testing.T) {
		_, err := bookRepository.Get(9999)
		require.EqualError(t, err, "some-error")
	})
}

func TestBookRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WillReturnRows(
			sqlmock.NewRows(bookColumns).
				AddRow(1, "title1", "author1", timekit.UTC("2019-02-20T10:00:00-05:00")).
				AddRow(2, "title2", "author2", timekit.UTC("2019-02-20T10:00:01-05:00")),
		)
	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("some-error"))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"wrong column"}).AddRow("data"))

	bookRepository := NewBookRepository(db)

	t.Run("return rows", func(t *testing.T) {
		books, err := bookRepository.List()
		require.NoError(t, err)
		require.Equal(t, books, []Book{
			{1, "title1", "author1", timekit.UTC("2019-02-20T10:00:00-05:00")},
			{2, "title2", "author2", timekit.UTC("2019-02-20T10:00:01-05:00")},
		})
	})

	t.Run("return error when database problem", func(t *testing.T) {
		_, err := bookRepository.List()
		require.EqualError(t, err, "some-error")
	})

	t.Run("return error when scan problem", func(t *testing.T) {
		_, err := bookRepository.List()
		require.EqualError(t, err, "sql: expected 1 destination arguments in Scan, not 4")
	})
}
