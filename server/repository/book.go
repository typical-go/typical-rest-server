package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"go.uber.org/dig"
)

// Book represented database model
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"created_at"`
}

// BookRepo to get book data from database [mock]
type BookRepo interface {
	FindOne(context.Context, int64) (*Book, error)
	Find(context.Context, ...dbkit.FindOption) ([]*Book, error)
	Create(context.Context, *Book) (*Book, error)
	Delete(context.Context, int64) error
	Update(context.Context, int64, *Book) (*Book, error)
}

// BookRepoImpl is implementation book repository
type BookRepoImpl struct {
	dig.In
	*typpostgres.DB
}

// NewBookRepo return new instance of BookRepo [constructor]
func NewBookRepo(impl BookRepoImpl) BookRepo {
	return &impl
}

// FindOne book
func (r *BookRepoImpl) FindOne(ctx context.Context, id int64) (*Book, error) {
	row := sq.
		Select(
			"id",
			"title",
			"author",
			"updated_at",
			"created_at",
		).
		From("books").
		Where(
			sq.Eq{"id": id},
		).
		PlaceholderFormat(sq.Dollar).
		RunWith(r).
		QueryRowContext(ctx)

	book := new(Book)
	if err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.UpdatedAt,
		&book.CreatedAt,
	); err != nil {
		return nil, err
	}

	return book, nil
}

// Find book
func (r *BookRepoImpl) Find(ctx context.Context, opts ...dbkit.FindOption) (list []*Book, err error) {
	var (
		rows *sql.Rows
	)
	builder := sq.
		Select(
			"id",
			"title",
			"author",
			"updated_at",
			"created_at",
		).
		From("books").
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	for _, opt := range opts {
		if builder, err = opt.CompileQuery(builder); err != nil {
			return
		}

	}
	if rows, err = builder.QueryContext(ctx); err != nil {
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
func (r *BookRepoImpl) Create(ctx context.Context, book *Book) (*Book, error) {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	scanner := sq.
		Insert("books").
		Columns(
			"title",
			"author",
			"created_at",
			"updated_at",
		).
		Values(
			book.Title,
			book.Author,
			book.CreatedAt,
			book.UpdatedAt,
		).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r).
		QueryRowContext(ctx)

	if err := scanner.Scan(&book.ID); err != nil {
		return nil, err
	}
	return book, nil
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	query := sq.
		Delete("books").
		Where(
			sq.Eq{"id": id},
		).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if _, err = query.ExecContext(ctx); err != nil {
		return
	}
	return
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, id int64, forms *Book) (book *Book, err error) {
	if book, err = r.FindOne(ctx, id); err != nil {
		return
	}

	book.Title = forms.Title
	book.Author = forms.Author
	book.UpdatedAt = time.Now()

	update := sq.Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("updated_at", book.UpdatedAt).
		Where(
			sq.Eq{"id": id},
		).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	if _, err = update.ExecContext(ctx); err != nil {
		return nil, err
	}
	return
}
