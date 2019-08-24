package repository

import (
	"context"
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
	Find(ctx context.Context, id int64) (*Book, error)
	List(ctx context.Context) ([]*Book, error)
	Insert(ctx context.Context, book Book) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, book Book) error
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
func (r *BookRepositoryImpl) Find(ctx context.Context, id int64) (book *Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		Where(sq.Eq{"id": id})

	rows, err := builder.RunWith(r.DB).QueryContext(ctx)
	if err != nil {
		return
	}

	if rows.Next() {
		book, err = scanBook(rows)
	}
	return
}

// List book
func (r *BookRepositoryImpl) List(ctx context.Context) (list []*Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "title", "author", "updated_at", "created_at").From("books")

	rows, err := builder.RunWith(r.DB).QueryContext(ctx)
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
func (r *BookRepositoryImpl) Insert(ctx context.Context, book Book) (lastInsertID int64, err error) {
	query := sq.Insert("books").
		Columns("title", "author").
		Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRowContext(ctx).Scan(&book.ID)
	if err != nil {
		return
	}

	lastInsertID = book.ID
	return
}

// Delete book
func (r *BookRepositoryImpl) Delete(ctx context.Context, id int64) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("books").
		Where(sq.Eq{"id": id})

	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}

// Update book
func (r *BookRepositoryImpl) Update(ctx context.Context, book Book) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": book.ID})

	_, err = builder.RunWith(r.DB).ExecContext(ctx)
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
