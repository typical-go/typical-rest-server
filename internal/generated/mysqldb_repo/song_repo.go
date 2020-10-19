package mysqldb_repo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/reflectkit"
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
	// @mock
	SongRepo interface {
		Find(context.Context, ...dbkit.SelectOption) ([]*mysqldb.Song, error)
		Create(context.Context, *mysqldb.Song) (int64, error)
		Delete(context.Context, dbkit.DeleteOption) (int64, error)
		Update(context.Context, *mysqldb.Song, dbkit.UpdateOption) (int64, error)
		Patch(context.Context, *mysqldb.Song, dbkit.UpdateOption) (int64, error)
	}
	// SongRepoImpl is implementation songs repository
	SongRepoImpl struct {
		dig.In
		*sql.DB `name:"mysql"`
	}
)

func init() {
	typapp.AppendCtor(&typapp.Constructor{Name: "", Fn: NewSongRepo})
}

// NewSongRepo return new instance of SongRepo
func NewSongRepo(impl SongRepoImpl) SongRepo {
	return &impl
}

// Find songs
func (r *SongRepoImpl) Find(ctx context.Context, opts ...dbkit.SelectOption) (list []*mysqldb.Song, err error) {
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

	list = make([]*mysqldb.Song, 0)
	for rows.Next() {
		ent := new(mysqldb.Song)
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

// Create songs
func (r *SongRepoImpl) Create(ctx context.Context, ent *mysqldb.Song) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	res, err := sq.
		Insert(SongTableName).
		Columns(
			SongTable.Title,
			SongTable.Artist,
			SongTable.UpdatedAt,
			SongTable.CreatedAt,
		).
		Values(
			ent.Title,
			ent.Artist,
			time.Now(),
			time.Now(),
		).
		RunWith(txn.DB).
		ExecContext(ctx)

	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	lastInsertID, err := res.LastInsertId()
	txn.SetError(err)
	return lastInsertID, err
}

// Update songs
func (r *SongRepoImpl) Update(ctx context.Context, ent *mysqldb.Song, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(SongTableName).
		Set(SongTable.Title, ent.Title).
		Set(SongTable.Artist, ent.Artist).
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
	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Patch songs
func (r *SongRepoImpl) Patch(ctx context.Context, ent *mysqldb.Song, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Update(SongTableName).RunWith(txn.DB)

	if !reflectkit.IsZero(ent.Title) {
		builder = builder.Set(SongTable.Title, ent.Title)
	}
	if !reflectkit.IsZero(ent.Artist) {
		builder = builder.Set(SongTable.Artist, ent.Artist)
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

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Delete songs
func (r *SongRepoImpl) Delete(ctx context.Context, opt dbkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.Delete(SongTableName).RunWith(txn.DB)
	if builder, err = opt.CompileDelete(builder); err != nil {
		txn.SetError(err)
		return -1, err
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
