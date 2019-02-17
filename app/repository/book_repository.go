package repository

import (
	"database/sql"
	"fmt"
	"time"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

// Book represented database model
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// BookRepository to get book data from databasesa
type BookRepository interface {
	Get(id int) (*Book, error)
	List() ([]Book, error)
	Insert(book Book) error
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

func (r *bookRepository) Get(id int) (book *Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("id", "title", "author", "created_at").
		From("books").Where(sq.Eq{"id": id})

	rows, err := builder.RunWith(r.conn).Query()
	if err != nil {
		return
	}

	if rows.Next() {
		book = new(Book)
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt)
	}
	return
}

func (r *bookRepository) List() (list []Book, err error) {
	err = fmt.Errorf("Under Construction")
	return
}

func (r *bookRepository) Insert(book Book) (err error) {
	err = fmt.Errorf("Under Construction")
	return
}
