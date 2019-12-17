package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "title", "author", "updated_at", "created_at").
		From("books").
		Where(sq.Eq{"id": id})
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		var book0 Book
		if err = rows.Scan(&book0.ID, &book0.Title, &book0.Author, &book0.UpdatedAt, &book0.CreatedAt); err != nil {
			return nil, err
		}
		book = &book0
	}
	return
}

// List book
func (r *BookRepoImpl) List(ctx context.Context) (list []*Book, err error) {
	var rows *sql.Rows
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "title", "author", "updated_at", "created_at").
		From("books")
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	list = make([]*Book, 0)
	for rows.Next() {
		var book0 Book
		if err = rows.Scan(&book0.ID, &book0.Title, &book0.Author, &book0.UpdatedAt, &book0.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &book0)
	}
	return
}

// Insert book
func (r *BookRepoImpl) Insert(ctx context.Context, book Book) (lastInsertID int64, err error) {
	query := sq.Insert("books").
		Columns("title", "author").
		Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)
	if err = query.QueryRowContext(ctx).Scan(&book.ID); err != nil {
		return
	}
	lastInsertID = book.ID
	return
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("books").Where(sq.Eq{"id": id})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, book Book) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": book.ID})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}
