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

// BookRepoImpl is implementation book repository
type BookRepoImpl struct {
	dig.In
	*typpostgres.DB
}

// FindOne book
func (r *BookRepoImpl) FindOne(ctx context.Context, id int64) (*Book, error) {
	builder := sq.
		Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)

	e := new(Book)
	if err := builder.QueryRowContext(ctx).Scan(&e.ID, &e.Title, &e.Author, &e.UpdatedAt, &e.CreatedAt); err != nil {
		return nil, err
	}

	return e, nil
}

// Find book
func (r *BookRepoImpl) Find(ctx context.Context, opts ...dbkit.FindOption) (list []*Book, err error) {
	var (
		rows *sql.Rows
	)
	builder := sq.
		Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		PlaceholderFormat(sq.Dollar).RunWith(r)

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
		var book0 Book
		if err = rows.Scan(&book0.ID, &book0.Title, &book0.Author, &book0.UpdatedAt, &book0.CreatedAt); err != nil {
			return
		}
		list = append(list, &book0)
	}
	return
}

// Create book
func (r *BookRepoImpl) Create(ctx context.Context, new *Book) (book *Book, err error) {
	book = &Book{
		Title:     new.Title,
		Author:    new.Author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	builder := sq.
		Insert("books").
		Columns("title", "author", "created_at", "updated_at").
		Values(book.Title, book.Author, book.CreatedAt, book.UpdatedAt).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if err = builder.QueryRowContext(ctx).Scan(&book.ID); err != nil {
		return
	}
	return
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	query := sq.
		Delete("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if _, err = query.ExecContext(ctx); err != nil {
		return
	}
	return
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, update *Book) (book *Book, err error) {
	if book, err = r.FindOne(ctx, update.ID); err != nil {
		return
	}
	book.Title = update.Title
	book.Author = update.Author
	book.UpdatedAt = time.Now()
	builder := sq.Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("updated_at", book.UpdatedAt).
		Where(sq.Eq{"id": book.ID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)
	if _, err = builder.ExecContext(ctx); err != nil {
		return nil, err
	}
	return
}
