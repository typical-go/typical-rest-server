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
		music, err = scanMusic(rows)
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
		var music *Music
		if music, err = scanMusic(rows); err != nil {
			return
		}
		list = append(list, music)
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

func scanMusic(rows *sql.Rows) (*Music, error) {
	var music Music
	var err error
	if err = rows.Scan(&music.ID, &music.Artist, &music.UpdatedAt, &music.CreatedAt); err != nil {
		return nil, err
	}
	return &music, nil
}
