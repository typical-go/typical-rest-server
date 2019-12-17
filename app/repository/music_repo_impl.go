package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "artist", "updated_at", "created_at").
		From("musics").
		Where(sq.Eq{"id": id})
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		var music0 Music
		if err = rows.Scan(&music0.ID, &music0.Artist, &music0.UpdatedAt, &music0.CreatedAt); err != nil {
			return nil, err
		}
	}
	return
}

// List music
func (r *MusicRepoImpl) List(ctx context.Context) (list []*Music, err error) {
	var rows *sql.Rows
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("id", "artist", "updated_at", "created_at").
		From("musics")
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	list = make([]*Music, 0)
	for rows.Next() {
		var music0 Music
		if err = rows.Scan(&music0.ID, &music0.Artist, &music0.UpdatedAt, &music0.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &music0)
	}
	return
}

// Insert music
func (r *MusicRepoImpl) Insert(ctx context.Context, music Music) (lastInsertID int64, err error) {
	query := sq.Insert("musics").
		Columns("artist").
		Values(music.Artist).
		Suffix("RETURNING \"id\"").
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)
	if err = query.QueryRowContext(ctx).Scan(&music.ID); err != nil {
		return
	}
	lastInsertID = music.ID
	return
}

// Delete music
func (r *MusicRepoImpl) Delete(ctx context.Context, id int64) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("musics").Where(sq.Eq{"id": id})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}

// Update music
func (r *MusicRepoImpl) Update(ctx context.Context, music Music) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("musics").
		Set("artist", music.Artist).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": music.ID})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}
