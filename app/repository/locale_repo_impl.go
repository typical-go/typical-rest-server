package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"go.uber.org/dig"
)

// LocaleRepoImpl is implementation locale repository
type LocaleRepoImpl struct {
	dig.In
	*sql.DB
}

// Find locale
func (r *LocaleRepoImpl) Find(ctx context.Context, id int64) (e *Locale, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "lang_code", "country_code", "updated_at", "created_at").
		From("locales").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		e = new(Locale)
		if err = rows.Scan(&e.ID, &e.LangCode, &e.CountryCode, &e.UpdatedAt, &e.CreatedAt); err != nil {
			return nil, err
		}
	}
	return
}

// List locale
func (r *LocaleRepoImpl) List(ctx context.Context) (list []*Locale, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "lang_code", "country_code", "updated_at", "created_at").
		From("locales").
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	list = make([]*Locale, 0)
	for rows.Next() {
		var e0 Locale
		if err = rows.Scan(&e0.ID, &e0.LangCode, &e0.CountryCode, &e0.UpdatedAt, &e0.CreatedAt); err != nil {
			return
		}
		list = append(list, &e0)
	}
	return
}

// Insert locale
func (r *LocaleRepoImpl) Insert(ctx context.Context, e Locale) (lastInsertID int64, err error) {
	builder := sq.
		Insert("locales").
		Columns("lang_code", "country_code").
		Values(e.LangCode, e.CountryCode).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&e.ID); err != nil {
		return
	}
	lastInsertID = e.ID
	return
}

// Delete locale
func (r *LocaleRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	builder := sq.
		Delete("locales").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}

// Update locale
func (r *LocaleRepoImpl) Update(ctx context.Context, e Locale) (err error) {
	builder := sq.
		Update("locales").
		Set("lang_code", e.LangCode).
		Set("country_code", e.CountryCode).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": e.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
