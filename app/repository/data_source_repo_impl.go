package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"go.uber.org/dig"
)

// DataSourceRepoImpl is implementation data_source repository
type DataSourceRepoImpl struct {
	dig.In
	*sql.DB
}

// Find data_source
func (r *DataSourceRepoImpl) Find(ctx context.Context, id int64) (e *DataSource, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "name", "url", "updated_at", "created_at").
		From("data_sources").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		e = new(DataSource)
		if err = rows.Scan(&e.ID, &e.Name, &e.Url, &e.UpdatedAt, &e.CreatedAt); err != nil {
			return nil, err
		}
	}
	return
}

// List data_source
func (r *DataSourceRepoImpl) List(ctx context.Context) (list []*DataSource, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "name", "url", "updated_at", "created_at").
		From("data_sources").
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	list = make([]*DataSource, 0)
	for rows.Next() {
		var e0 DataSource
		if err = rows.Scan(&e0.ID, &e0.Name, &e0.Url, &e0.UpdatedAt, &e0.CreatedAt); err != nil {
			return
		}
		list = append(list, &e0)
	}
	return
}

// Insert data_source
func (r *DataSourceRepoImpl) Insert(ctx context.Context, e DataSource) (lastInsertID int64, err error) {
	builder := sq.
		Insert("data_sources").
		Columns("name", "url").
		Values(e.Name, e.Url).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&e.ID); err != nil {
		return
	}
	lastInsertID = e.ID
	return
}

// Delete data_source
func (r *DataSourceRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	builder := sq.
		Delete("data_sources").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}

// Update data_source
func (r *DataSourceRepoImpl) Update(ctx context.Context, e DataSource) (err error) {
	builder := sq.
		Update("data_sources").
		Set("name", e.Name).
		Set("url", e.Url).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": e.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
