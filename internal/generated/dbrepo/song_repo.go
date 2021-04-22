package dbrepo

/* DO NOT EDIT. This file generated due to '@dbrepo' annotation */

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-rest-server/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/reflectkit"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"
)

var (
	// SongTableName is table name for songs entity
	SongTableName = "songs"
	// SongTable is columns for songs entity
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
	// SongRepo to get songs data from database
	SongRepo interface {
		Count(context.Context, ...sqkit.SelectOption) (int64, error)
		Find(context.Context, ...sqkit.SelectOption) ([]*entity.Song, error)
		Insert(context.Context, ...*entity.Song) (int64, error)
		Delete(context.Context, sqkit.DeleteOption) (int64, error)
		Update(context.Context, *entity.Song, sqkit.UpdateOption) (int64, error)
		Patch(context.Context, *entity.Song, sqkit.UpdateOption) (int64, error)
	}
	// SongRepoImpl is implementation songs repository
	SongRepoImpl struct {
		dig.In
		*sql.DB `name:"mysql"`
	}
)

func init() {
	typapp.Provide("", NewSongRepo)
}

// NewSongRepo return new instance of SongRepo
func NewSongRepo(impl SongRepoImpl) SongRepo {
	return &impl
}

// Count songs
func (r *SongRepoImpl) Count(ctx context.Context, opts ...sqkit.SelectOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}
	builder := sq.
		Select("count(*)").
		From(SongTableName).
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

// Find songs
func (r *SongRepoImpl) Find(ctx context.Context, opts ...sqkit.SelectOption) (list []*entity.Song, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	builder := sq.
		Select(
			SongTable.ID,
			SongTable.Title,
			SongTable.Artist,
			SongTable.UpdatedAt,
			SongTable.CreatedAt,
		).
		From(SongTableName).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*entity.Song, 0)
	for rows.Next() {
		ent := new(entity.Song)
		if err = rows.Scan(
			&ent.ID,
			&ent.Title,
			&ent.Artist,
			&ent.UpdatedAt,
			&ent.CreatedAt,
		); err != nil {
			return
		}
		list = append(list, ent)
	}
	return
}

// Insert songs
func (r *SongRepoImpl) Insert(ctx context.Context, ents ...*entity.Song) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Insert(SongTableName).
		Columns(
			SongTable.Title,
			SongTable.Artist,
			SongTable.UpdatedAt,
			SongTable.CreatedAt,
		)

	for _, ent := range ents {
		builder = builder.Values(
			ent.Title,
			ent.Artist,
			time.Now(),
			time.Now(),
		)
	}

	res, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	lastInsertID, err := res.LastInsertId()
	txn.SetError(err)
	return lastInsertID, err
}

// Update songs
func (r *SongRepoImpl) Update(ctx context.Context, ent *entity.Song, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(SongTableName).
		Set(SongTable.Title, ent.Title).
		Set(SongTable.Artist, ent.Artist).
		Set(SongTable.UpdatedAt, time.Now()).
		RunWith(txn)

	if opt != nil {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Patch songs
func (r *SongRepoImpl) Patch(ctx context.Context, ent *entity.Song, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Update(SongTableName).RunWith(txn)

	if !reflectkit.IsZero(ent.Title) {
		builder = builder.Set(SongTable.Title, ent.Title)
	}
	if !reflectkit.IsZero(ent.Artist) {
		builder = builder.Set(SongTable.Artist, ent.Artist)
	}
	builder = builder.Set(SongTable.UpdatedAt, time.Now())

	if opt != nil {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Delete songs
func (r *SongRepoImpl) Delete(ctx context.Context, opt sqkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Delete(SongTableName).RunWith(txn)
	if opt != nil {
		builder = opt.CompileDelete(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}
