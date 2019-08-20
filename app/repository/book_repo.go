package repository

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/dig"
)

// Book represented database model
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// BookRepository to get book data from databasesa
type BookRepository interface {
	Find(id int64) (*Book, error)
	List() ([]*Book, error)
	Insert(book Book) (lastInsertID int64, err error)
	Delete(id int64) error
	Update(book Book) error
}

// BookRepositoryImpl is implementation book repository
type BookRepositoryImpl struct {
	dig.In
	*sql.DB
}

// NewBookRepository return new instance of BookRepository
func NewBookRepository(impl BookRepositoryImpl) BookRepository {
	return &impl
}

// Find book
func (r *BookRepositoryImpl) Find(id int64) (book *Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).
		From(bookTable).
		Where(sq.Eq{idColumn: id})

	rows, err := builder.RunWith(r.DB).Query()
	if err != nil {
		return
	}

	if rows.Next() {
		book, err = scanBook(rows)
	}
	return
}

// List book
func (r *BookRepositoryImpl) List() (list []*Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).From(bookTable)

	rows, err := builder.RunWith(r.DB).Query()
	if err != nil {
		return
	}

	list = make([]*Book, 0)

	for rows.Next() {
		var book *Book
		book, err = scanBook(rows)
		if err != nil {
			return
		}
		list = append(list, book)
	}
	return
}

// Insert book
func (r *BookRepositoryImpl) Insert(book Book) (lastInsertID int64, err error) {
	query := sq.Insert(bookTable).
		Columns(bookTitleColumn, bookAuthorColumn).
		Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&book.ID)
	if err != nil {
		return
	}

	lastInsertID = book.ID
	return
}

// Delete book
func (r *BookRepositoryImpl) Delete(id int64) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete(bookTable).
		Where(sq.Eq{idColumn: id})

	_, err = builder.RunWith(r.DB).Exec()
	return
}

// Update book
func (r *BookRepositoryImpl) Update(book Book) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update(bookTable).
		Set(bookTitleColumn, book.Title).
		Set(bookAuthorColumn, book.Author).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: book.ID})

	_, err = builder.RunWith(r.DB).Exec()
	return
}

func scanBook(rows *sql.Rows) (*Book, error) {
	var book Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.UpdatedAt, &book.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &book, nil
}
