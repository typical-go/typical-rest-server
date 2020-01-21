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

// MusicRepoImpl is implementation music repository
type MusicRepoImpl struct {
	dig.In
	*typpostgres.DB
}

// FindOne music
func (r *MusicRepoImpl) FindOne(ctx context.Context, id int64) (e *Music, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "artist", "updated_at", "created_at").
		From("musics").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		e = new(Music)
		if err = rows.Scan(&e.ID, &e.Artist, &e.UpdatedAt, &e.CreatedAt); err != nil {
			return nil, err
		}
	}
	return
}

// Find musics
func (r *MusicRepoImpl) Find(ctx context.Context) (list []*Music, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "artist", "updated_at", "created_at").
		From("musics").
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	list = make([]*Music, 0)
	for rows.Next() {
		var e Music
		if err = rows.Scan(&e.ID, &e.Artist, &e.UpdatedAt, &e.CreatedAt); err != nil {
			return
		}
		list = append(list, &e)
	}
	return
}

// Create music
func (r *MusicRepoImpl) Create(ctx context.Context, e Music) (lastInsertID int64, err error) {
	builder := sq.
		Insert("musics").
		Columns("artist").
		Values(e.Artist).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&e.ID); err != nil {
		return
	}
	lastInsertID = e.ID
	return
}

// Delete music
func (r *MusicRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	builder := sq.
		Delete("musics").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}

// Update music
func (r *MusicRepoImpl) Update(ctx context.Context, e Music) (err error) {
	builder := sq.
		Update("musics").
		Set("artist", e.Artist).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": e.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(dbkit.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
