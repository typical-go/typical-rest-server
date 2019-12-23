package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"

	"go.uber.org/dig"
)

// BookRepoImpl is implementation book repository
type BookRepoImpl struct {
	dig.In
	*sql.DB
}

// Find book
func (r *BookRepoImpl) Find(ctx context.Context, id int64) (book *Book, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		var book0 Book
		if err = rows.Scan(&book0.ID, &book0.Title, &book0.Author, &book0.UpdatedAt, &book0.CreatedAt); err != nil {
			return
		}
		book = &book0
	}
	return
}

// List book
func (r *BookRepoImpl) List(ctx context.Context) (list []*Book, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		PlaceholderFormat(sq.Dollar).RunWith(r)
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

// Insert book
func (r *BookRepoImpl) Insert(ctx context.Context, book Book) (lastInsertID int64, err error) {
	builder := sq.
		Insert("books").
		Columns("title", "author").Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&book.ID); err != nil {
		return
	}
	lastInsertID = book.ID
	return
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	query := sq.
		Delete("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = query.ExecContext(ctx)
	return
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, book Book) (err error) {
	builder := sq.
		Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": book.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
