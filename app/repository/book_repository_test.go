package repository

import (
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/imantung/typical-go-server/config"
	"github.com/imantung/typical-go-server/db"
	"github.com/stretchr/testify/require"
)

func TestBookRepository(t *testing.T) {
	conf, _ := config.NewConfig()
	conf.DbName = fmt.Sprintf("%s_test", conf.DbName)

	// FIXME: migration outside the test
	// FIXME: clear database without drop
	err := db.ResetTestDB(conf, "file://../../db/migrate")
	require.NoError(t, err)

	conn, err := db.Connect(conf)
	require.NoError(t, err)
	defer conn.Close()

	repository := NewBookRepository(conn)

	var id int64
	var book *Book

	t.Run("insert new record", func(t *testing.T) {
		id, err = repository.Insert(Book{Title: "some-title", Author: "some-author"})
		require.NoError(t, err)
		require.True(t, id > 0)

		book, err = repository.Get(id)
		require.NoError(t, err)
		require.Equal(t, "some-title", book.Title)
		require.Equal(t, "some-author", book.Author)
	})

	t.Run("update record", func(t *testing.T) {
		err = repository.Update(Book{ID: id, Title: "new-title", Author: "new-author"})
		require.NoError(t, err)

		book, err = repository.Get(id)
		require.NoError(t, err)
		require.Equal(t, "new-title", book.Title)
		require.Equal(t, "new-author", book.Author)
	})

	t.Run("delete record", func(t *testing.T) {
		err = repository.Delete(id)
		require.NoError(t, err)

		book, err := repository.Get(id)
		require.NoError(t, err)
		require.Nil(t, book)
	})

	t.Run("list of record", func(t *testing.T) {
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
	})
}
