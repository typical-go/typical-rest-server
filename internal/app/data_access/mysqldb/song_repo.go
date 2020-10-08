package mysqldb

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
)

var (
	// SongTableName is table name for song entity
	SongTableName = "songs"
	// SongTable is columns for song entity
	SongTable = struct {
		ID        string
		Title     string
		Artist    string
		UpdatedAt string
		CreatedAt string
	}{
		ID:        "id",
		Title:     "title",
		Artist:    "artist",
		UpdatedAt: "updated_at",
		CreatedAt: "created_at",
	}
)

type (
	// SongRepo to get song data from database
	// @mock
	SongRepo interface {
		Find(context.Context, ...dbkit.SelectOption) ([]*Song, error)
		Create(context.Context, *Song) (int64, error)
		Delete(context.Context, dbkit.DeleteOption) (int64, error)
		Update(context.Context, *Song, dbkit.UpdateOption) (int64, error)
		Patch(context.Context, *Song, dbkit.UpdateOption) (int64, error)
	}
	// SongRepoImpl is implementation song repository
	SongRepoImpl struct {
		dig.In
		*sql.DB `name:"mysql"`
	}
)

// NewSongRepo return new instance of SongRepo
// @ctor
func NewSongRepo(impl SongRepoImpl) SongRepo {
	return &impl
}

// Find song
func (r *SongRepoImpl) Find(ctx context.Context, opts ...dbkit.SelectOption) (list []*Song, err error) {
	builder := sq.
		Select(
			SongTable.ID,
			SongTable.Title,
			SongTable.Artist,
			SongTable.UpdatedAt,
			SongTable.CreatedAt,
		).
		From(SongTableName).
		RunWith(r)

	for _, opt := range opts {
		if builder, err = opt.CompileSelect(builder); err != nil {
			return nil, err
		}
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*Song, 0)
	for rows.Next() {
		song := new(Song)
		if err = rows.Scan(
			&song.ID,
			&song.Title,
			&song.Artist,
			&song.UpdatedAt,
			&song.CreatedAt,
		); err != nil {
			return
		}
		list = append(list, song)
	}
	return
}

// Create song
func (r *SongRepoImpl) Create(ctx context.Context, song *Song) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	res, err := sq.
		Insert(SongTableName).
		Columns(
			SongTable.Title,
			SongTable.Artist,
			SongTable.CreatedAt,
			SongTable.UpdatedAt,
		).
		Values(
			song.Title,
			song.Artist,
			time.Now(),
			time.Now(),
		).
		RunWith(txn.DB).
		ExecContext(ctx)

	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	return res.LastInsertId()
}

// Delete song
func (r *SongRepoImpl) Delete(ctx context.Context, opt dbkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Delete(SongTableName).
		RunWith(txn.DB)

	if builder, err = opt.CompileDelete(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	return res.RowsAffected()
}

// Update song
func (r *SongRepoImpl) Update(ctx context.Context, song *Song, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(SongTableName).
		Set(SongTable.Title, song.Title).
		Set(SongTable.Artist, song.Artist).
		Set(SongTable.UpdatedAt, time.Now()).
		RunWith(txn.DB)

	if builder, err = opt.CompileUpdate(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	return res.RowsAffected()
}

// Patch song to update field of song if available
func (r *SongRepoImpl) Patch(ctx context.Context, song *Song, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(SongTableName).
		RunWith(txn.DB)

	if song.Title != "" {
		builder = builder.Set(SongTable.Title, song.Title)
	}
	if song.Artist != "" {
		builder = builder.Set(SongTable.Artist, song.Artist)
	}
	builder = builder.Set(SongTable.UpdatedAt, time.Now())

	if builder, err = opt.CompileUpdate(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	return res.RowsAffected()
}
