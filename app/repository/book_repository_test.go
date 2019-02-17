package repository_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/imantung/go-helper/timekit"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/stretchr/testify/require"
)

func TestBookRepository(t *testing.T) {
	expected := &repository.Book{
		ID:        1,
		Title:     "one",
		Author:    "author1",
		CreatedAt: timekit.UTC("2002-10-02T10:00:00-05:00"),
	}

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "author", "created_at"}).
		AddRow(expected.ID, expected.Title, expected.Author, expected.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	bookRepository := repository.NewBookRepository(db)

	book, err := bookRepository.Get(1)
	require.NoError(t, err)
	require.Equal(t, book, expected)
}
