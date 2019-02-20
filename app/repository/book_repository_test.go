package repository

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/imantung/go-helper/timekit"
	"github.com/stretchr/testify/require"
)

func TestBookRepository(t *testing.T) {
	expected := Book{
		ID:        1,
		Title:     "one",
		Author:    "author1",
		CreatedAt: timekit.UTC("2002-10-02T10:00:00-05:00"),
	}

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows(bookColumns).
		AddRow(expected.ID, expected.Title, expected.Author, expected.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)
	mock.ExpectQuery("SELECT").WithArgs(9999).WillReturnError(fmt.Errorf("some-error"))

	bookRepository := NewBookRepository(db)

	t.Run("return rows", func(t *testing.T) {
		book, err := bookRepository.Get(1)
		require.NoError(t, err)
		require.Equal(t, book, expected)
	})

	t.Run("return error", func(t *testing.T) {
		_, err := bookRepository.Get(9999)
		require.EqualError(t, err, "some-error")
	})

}
