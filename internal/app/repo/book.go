package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/reflectkit"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"

	sq "github.com/Masterminds/squirrel"
)

// TODO: create repo code generator from database table

type (
	Book struct {
		ID        int64     `column:"id" option:"pk" json:"id"`
		Title     string    `column:"title" json:"title"`
		Author    string    `column:"author" json:"author"`
		UpdatedAt time.Time `column:"updated_at" option:"now" json:"update_at"`
		CreatedAt time.Time `column:"created_at" option:"now,no_update" json:"created_at"`
	}
	// @mock
	BookRepo interface {
		Count(context.Context, ...sqkit.SelectOption) (int64, error)
		Find(context.Context, ...sqkit.SelectOption) ([]*Book, error)
		Insert(context.Context, *Book) (int64, error)
		BulkInsert(context.Context, ...*Book) (int64, error)
		Delete(context.Context, ...sqkit.DeleteOption) (int64, error)
		Update(context.Context, *Book, ...sqkit.UpdateOption) (int64, error)
		Patch(context.Context, *Book, ...sqkit.UpdateOption) (int64, error)
	}
	BookRepoImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}
)

var (
	BookTableName = "books"
	BookTable     = struct {
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
)

// @ctor
func NewBookRepo(impl BookRepoImpl) BookRepo {
	return &impl
}

// Count books
func (r *BookRepoImpl) Count(ctx context.Context, opts ...sqkit.SelectOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}
	builder := sq.
		Select("count(1)").
		From(BookTableName).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	row := builder.QueryRowContext(ctx)

	var cnt int64
	if err := row.Scan(&cnt); err != nil {
		return -1, err
	}
	return cnt, nil
}

// Find books
func (r *BookRepoImpl) Find(ctx context.Context, opts ...sqkit.SelectOption) (list []*Book, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	builder := sq.
		Select(
			BookTable.ID,
			BookTable.Title,
			BookTable.Author,
			BookTable.UpdatedAt,
			BookTable.CreatedAt,
		).
		From(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*Book, 0)
	for rows.Next() {
		ent := new(Book)
		if err = rows.Scan(
			&ent.ID,
			&ent.Title,
			&ent.Author,
			&ent.UpdatedAt,
			&ent.CreatedAt,
		); err != nil {
			return
		}
		list = append(list, ent)
	}
	return
}

// Insert books and return last inserted id
func (r *BookRepoImpl) Insert(ctx context.Context, ent *Book) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Insert(BookTableName).
		Columns(
			BookTable.Title,
			BookTable.Author,
			BookTable.UpdatedAt,
			BookTable.CreatedAt,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", BookTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		Values(
			ent.Title,
			ent.Author,
			time.Now(),
			time.Now(),
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		txn.AppendError(err)
		return -1, err
	}
	return id, nil
}

// BulkInsert books and return affected rows
func (r *BookRepoImpl) BulkInsert(ctx context.Context, ents ...*Book) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Insert(BookTableName).
		Columns(
			BookTable.Title,
			BookTable.Author,
			BookTable.UpdatedAt,
			BookTable.CreatedAt,
		).
		PlaceholderFormat(sq.Dollar)

	for _, ent := range ents {
		builder = builder.Values(
			ent.Title,
			ent.Author,
			time.Now(),
			time.Now(),
		)
	}

	res, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}

// Update books
func (r *BookRepoImpl) Update(ctx context.Context, ent *Book, opts ...sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(BookTableName).
		Set(BookTable.Title, ent.Title).
		Set(BookTable.Author, ent.Author).
		Set(BookTable.UpdatedAt, time.Now()).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}

// Patch books
func (r *BookRepoImpl) Patch(ctx context.Context, ent *Book, opts ...sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	if !reflectkit.IsZero(ent.Title) {
		builder = builder.Set(BookTable.Title, ent.Title)
	}
	if !reflectkit.IsZero(ent.Author) {
		builder = builder.Set(BookTable.Author, ent.Author)
	}
	builder = builder.Set(BookTable.UpdatedAt, time.Now())

	for _, opt := range opts {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}

// Delete books
func (r *BookRepoImpl) Delete(ctx context.Context, opts ...sqkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Delete(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileDelete(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}
