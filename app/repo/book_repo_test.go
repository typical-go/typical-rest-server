package repo

import (
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/imantung/typical-go-server/config"
	"github.com/stretchr/testify/require"
)

func TestBookRepository(t *testing.T) {
	conf, _ := config.LoadConfig()
	conn, err := conf.Postgres.OpenDBTest()
	require.NoError(t, err)
	defer conn.Close()

	repository := NewBookRepository(conn)

	var id int64
	var book *Book

	t.Run("Insert new record", func(t *testing.T) {
		id, err = repository.Insert(Book{Title: "some-title", Author: "some-author"})
		require.NoError(t, err)
		require.True(t, id > 0)

		book, err = repository.Get(id)
		require.NoError(t, err)
		require.Equal(t, "some-title", book.Title)
		require.Equal(t, "some-author", book.Author)
	})

	t.Run("Update record", func(t *testing.T) {
		err = repository.Update(Book{ID: id, Title: "new-title", Author: "new-author"})
		require.NoError(t, err)

		book, err = repository.Get(id)
		require.NoError(t, err)
		require.Equal(t, "new-title", book.Title)
		require.Equal(t, "new-author", book.Author)
	})

	t.Run("Delete record", func(t *testing.T) {
		err = repository.Delete(id)
		require.NoError(t, err)

		book, err := repository.Get(id)
		require.NoError(t, err)
		require.Nil(t, book)
	})
}

func TestBookRepository_List(t *testing.T) {
	conf, _ := config.LoadConfig()
	conn, err := conf.Postgres.Open()
	require.NoError(t, err)
	defer conn.Close()

	_, err = conn.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY`, bookTable))
	require.NoError(t, err)

	repository := NewBookRepository(conn)

	data := []Book{
		Book{Title: "title-1001", Author: "author-1001"},
		Book{Title: "title-1002", Author: "author-1002"},
		Book{Title: "title-1003", Author: "author-1003"},
	}

	for _, book := range data {
		repository.Insert(book)
	}

	list, err := repository.List()
	require.NoError(t, err)
	for i := range list {
		require.Equal(t, data[i].Title, list[i].Title)
		require.Equal(t, data[i].Author, list[i].Author)
	}

}
