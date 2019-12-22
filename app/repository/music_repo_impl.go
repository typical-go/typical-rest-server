package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"go.uber.org/dig"
)

// MusicRepoImpl is implementation music repository
type MusicRepoImpl struct {
	dig.In
	*sql.DB
}

// Find music
func (r *MusicRepoImpl) Find(ctx context.Context, id int64) (music *Music, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "artist", "updated_at", "created_at").
		From("musics").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		music = new(Music)
		if err = rows.Scan(&music.ID, &music.Artist, &music.UpdatedAt, &music.CreatedAt); err != nil {
			return nil, err
		}
	}
	return
}

// List music
func (r *MusicRepoImpl) List(ctx context.Context) (list []*Music, err error) {
	var rows *sql.Rows
	builder := sq.
		Select("id", "artist", "updated_at", "created_at").
		From("musics").
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	list = make([]*Music, 0)
	for rows.Next() {
		var music0 Music
		if err = rows.Scan(&music0.ID, &music0.Artist, &music0.UpdatedAt, &music0.CreatedAt); err != nil {
			return
		}
		list = append(list, &music0)
	}
	return
}

// Insert music
func (r *MusicRepoImpl) Insert(ctx context.Context, music Music) (lastInsertID int64, err error) {
	builder := sq.
		Insert("musics").
		Columns("artist").
		Values(music.Artist).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&music.ID); err != nil {
		return
	}
	lastInsertID = music.ID
	return
}

// Delete music
func (r *MusicRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	builder := sq.
		Delete("musics").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}

// Update music
func (r *MusicRepoImpl) Update(ctx context.Context, music Music) (err error) {
	builder := sq.
		Update("musics").
		Set("artist", music.Artist).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": music.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
