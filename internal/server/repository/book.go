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
		Retrieve(context.Context, ...dbkit.SelectOption) ([]*Book, error)
		Create(context.Context, *Book) (int64, error)
		Delete(context.Context, dbkit.DeleteOption) (int64, error)
		Update(context.Context, *Book, dbkit.UpdateOption) (int64, error)
		Patch(context.Context, *Book, dbkit.UpdateOption) (int64, error)
	}

	// BookRepoImpl is implementation book repository
	BookRepoImpl struct {
		dig.In
		*sql.DB
	}
)

// NewBookRepo return new instance of BookRepo
// @ctor
func NewBookRepo(impl BookRepoImpl) BookRepo {
	return &impl
}

// Retrieve book
func (r *BookRepoImpl) Retrieve(ctx context.Context, opts ...dbkit.SelectOption) (list []*Book, err error) {
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
			time.Now(),
			time.Now(),
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
func (r *BookRepoImpl) Delete(ctx context.Context, opt dbkit.DeleteOption) (int64, error) {
	builder := sq.
		Delete(BookTable).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	builder, err := opt.CompileDelete(builder)
	if err != nil {
		return -1, err
	}

	result, err := builder.ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, book *Book, opt dbkit.UpdateOption) (int64, error) {
	builder := sq.
		Update(BookTable).
		Set(BookCols.Title, book.Title).
		Set(BookCols.Author, book.Author).
		Set(BookCols.UpdatedAt, time.Now()).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	builder, err := opt.CompileUpdate(builder)
	if err != nil {
		return -1, err
	}

	result, err := builder.ExecContext(ctx)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Patch book to update field of book if available
func (r *BookRepoImpl) Patch(ctx context.Context, book *Book, opt dbkit.UpdateOption) (int64, error) {
	builder := sq.
		Update(BookTable).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	if book.Title != "" {
		builder = builder.Set(BookCols.Title, book.Title)
	}

	if book.Author != "" {
		builder = builder.Set(BookCols.Author, book.Author)
	}

	builder = builder.Set(BookCols.UpdatedAt, time.Now())

	builder, err := opt.CompileUpdate(builder)
	if err != nil {
		return -1, err
	}

	result, err := builder.ExecContext(ctx)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
