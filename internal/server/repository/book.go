package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"go.uber.org/dig"
)

var (
	// BookCols is columns for book entity
	BookCols = struct {
		ID        string
		Title     string
		Author    string
		UpdatedAt string
		CreatedAt string
	}{
		ID:        "id",
		Title:     "title",
		Author:    "author",
		UpdatedAt: "updated_at",
		CreatedAt: "created_at",
	}

	// BookTable is table name for book entity
	BookTable = "books"
)

type (
	// Book represented database model
	Book struct {
		ID        int64     `json:"id"`
		Title     string    `json:"title" validate:"required"`
		Author    string    `json:"author" validate:"required"`
		UpdatedAt time.Time `json:"update_at"`
		CreatedAt time.Time `json:"created_at"`
	}

	// BookRepo to get book data from database
	// @mock
	BookRepo interface {
		Find(context.Context, ...dbkit.SelectOption) ([]*Book, error)
		Create(context.Context, *Book) (int64, error)
		Delete(context.Context, dbkit.DeleteOption) error
		Update(context.Context, *Book, dbkit.UpdateOption) error
	}

	// BookRepoImpl is implementation book repository
	BookRepoImpl struct {
		dig.In
		*sql.DB
	}
)

// NewBookRepo return new instance of BookRepo
// @constructor
func NewBookRepo(impl BookRepoImpl) BookRepo {
	return &impl
}

// Find book
func (r *BookRepoImpl) Find(ctx context.Context, opts ...dbkit.SelectOption) (list []*Book, err error) {
	builder := sq.
		Select(
			BookCols.ID,
			BookCols.Title,
			BookCols.Author,
			BookCols.UpdatedAt,
			BookCols.CreatedAt,
		).
		From(BookTable).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	for _, opt := range opts {
		if builder, err = opt.CompileSelect(builder); err != nil {
			return nil, fmt.Errorf("book-repo: %w", err)
		}
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*Book, 0)
	for rows.Next() {
		book := new(Book)
		if err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.UpdatedAt,
			&book.CreatedAt,
		); err != nil {
			return
		}
		list = append(list, book)
	}
	return
}

// Create book
func (r *BookRepoImpl) Create(ctx context.Context, book *Book) (int64, error) {
	var id int64

	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	scanner := sq.
		Insert(BookTable).
		Columns(
			BookCols.Title,
			BookCols.Author,
			BookCols.CreatedAt,
			BookCols.UpdatedAt,
		).
		Values(
			book.Title,
			book.Author,
			book.CreatedAt,
			book.UpdatedAt,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", BookCols.ID),
		).
		PlaceholderFormat(sq.Dollar).
		RunWith(r).
		QueryRowContext(ctx)

	if err := scanner.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, opt dbkit.DeleteOption) (err error) {
	builder := sq.
		Delete(BookTable).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	if builder, err = opt.CompileDelete(builder); err != nil {
		return
	}

	if _, err = builder.ExecContext(ctx); err != nil {
		return
	}
	return
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, book *Book, opt dbkit.UpdateOption) (err error) {
	builder := sq.
		Update(BookTable).
		Set(BookCols.Title, book.Title).
		Set(BookCols.Author, book.Author).
		Set(BookCols.UpdatedAt, time.Now()).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	if builder, err = opt.CompileUpdate(builder); err != nil {
		return
	}

	_, err = builder.ExecContext(ctx)
	return
}
