package repository

import (
	"database/sql"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

// BookRepository to get book data from databasesa
type BookRepository interface {
	Get(id int) (Book, error)
	List() ([]Book, error)
	Insert(book Book) (sql.Result, error)
	Delete(id int) (sql.Result, error)
}

type bookRepository struct {
	conn *sql.DB
}

// NewBookRepository return new instance of BookRepository
func NewBookRepository(conn *sql.DB) BookRepository {
	return &bookRepository{
		conn: conn,
	}
}

func (r *bookRepository) Get(id int) (book Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(bookColumns...).
		From(bookTable).
		Where(sq.Eq{"id": id})

	rows, err := builder.RunWith(r.conn).Query()
	if err != nil {
		return
	}

	if rows.Next() {
		book, err = scanBook(rows)
	}
	return
}

func (r *bookRepository) List() (list []Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(bookColumns...).From(bookTable)

	rows, err := builder.RunWith(r.conn).Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var book Book
		book, err = scanBook(rows)
		if err != nil {
			return
		}
		list = append(list, book)
	}
	return
}

func (r *bookRepository) Insert(book Book) (result sql.Result, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Insert(bookTable).Columns("title", "author").
		Values(book.Author, book.Title)
	return builder.RunWith(r.conn).Exec()
}

func (r *bookRepository) Delete(id int) (result sql.Result, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete(bookTable).Where(sq.Eq{"id": id})
	return builder.RunWith(r.conn).Exec()
}
