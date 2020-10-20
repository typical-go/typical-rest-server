package service

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/internal/generated/mysqldb_repo"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// SongSvc contain logic for Song Controller
	// @mock
	SongSvc interface {
		FindOne(context.Context, string) (*mysqldb.Song, error)
		Find(context.Context) ([]*mysqldb.Song, error)
		Create(context.Context, *mysqldb.Song) (*mysqldb.Song, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *mysqldb.Song) (*mysqldb.Song, error)
		Patch(context.Context, string, *mysqldb.Song) (*mysqldb.Song, error)
	}
	// SongSvcImpl is implementation of SongSvc
	SongSvcImpl struct {
		dig.In
		mysqldb_repo.SongRepo
	}
)

// NewSongSvc return new instance of SongSvc
// @ctor
func NewSongSvc(impl SongSvcImpl) SongSvc {
	return &impl
}

// Create Song
func (b *SongSvcImpl) Create(ctx context.Context, book *mysqldb.Song) (*mysqldb.Song, error) {
	if err := validator.New().Struct(book); err != nil {
		return nil, typrest.NewValidErr(err.Error())
	}
	id, err := b.SongRepo.Create(ctx, book)
	if err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

// Find books
func (b *SongSvcImpl) Find(ctx context.Context) ([]*mysqldb.Song, error) {
	return b.SongRepo.Find(ctx)
}

// FindOne book
func (b *SongSvcImpl) FindOne(ctx context.Context, paramID string) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, typrest.NewValidErr("paramID is missing")
	}
	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) findOne(ctx context.Context, id int64) (*mysqldb.Song, error) {
	books, err := b.SongRepo.Find(ctx, dbkit.Equal(mysqldb_repo.SongTable.ID, id))
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
		return nil, sql.ErrNoRows
	}
	return books[0], nil
}

// Delete book
func (b *SongSvcImpl) Delete(ctx context.Context, paramID string) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return typrest.NewValidErr("paramID is missing")
	}
	_, err := b.SongRepo.Delete(ctx, dbkit.Equal(mysqldb_repo.SongTable.ID, id))
	return err
}

// Update book
func (b *SongSvcImpl) Update(ctx context.Context, paramID string, book *mysqldb.Song) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, typrest.NewValidErr("paramID is missing")
	}
	err := validator.New().Struct(book)
	if err != nil {
		return nil, typrest.NewValidErr(err.Error())
	}
	affectedRow, err := b.SongRepo.Update(ctx, book, dbkit.Equal(mysqldb_repo.SongTable.ID, id))
	if err != nil {
		return nil, err
	}
	if affectedRow < 1 {
		return nil, sql.ErrNoRows
	}
	return b.findOne(ctx, id)
}

// Patch book
func (b *SongSvcImpl) Patch(ctx context.Context, paramID string, book *mysqldb.Song) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, typrest.NewValidErr("paramID is missing")
	}
	affectedRow, err := b.SongRepo.Patch(ctx, book, dbkit.Equal(mysqldb_repo.SongTable.ID, id))
	if err != nil {
		return nil, err
	}
	if affectedRow < 1 {
		return nil, sql.ErrNoRows
	}
	return b.findOne(ctx, id)
}
