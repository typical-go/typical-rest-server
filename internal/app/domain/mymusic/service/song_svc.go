package service

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/internal/generated/mysqldb_repo"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// SongSvc contain logic for Song Controller
	// @mock
	SongSvc interface {
		FindOne(context.Context, string) (*mysqldb.Song, error)
		Find(context.Context, *FindReq) ([]*mysqldb.Song, error)
		Create(context.Context, *mysqldb.Song) (*mysqldb.Song, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *mysqldb.Song) (*mysqldb.Song, error)
		Patch(context.Context, string, *mysqldb.Song) (*mysqldb.Song, error)
	}
	// SongSvcImpl is implementation of SongSvc
	SongSvcImpl struct {
		dig.In
		Repo mysqldb_repo.SongRepo
	}
	// FindReq find request
	FindReq struct {
		Limit  uint64 `query:"limit"`
		Offset uint64 `query:"offset"`
		Sort   string `query:"sort"`
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
		return nil, echokit.NewValidErr(err.Error())
	}
	id, err := b.Repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

// Find books
func (b *SongSvcImpl) Find(ctx context.Context, req *FindReq) ([]*mysqldb.Song, error) {
	return b.Repo.Find(ctx, b.findSelectOpt(req)...)
}

func (b *SongSvcImpl) findSelectOpt(req *FindReq) (opts []dbkit.SelectOption) {
	opts = append(opts, &dbkit.OffsetPagination{Offset: req.Offset, Limit: req.Limit})
	if req.Sort != "" {
		opts = append(opts, dbkit.Sorts(strings.Split(req.Sort, ",")))
	}
	return
}

// FindOne book
func (b *SongSvcImpl) FindOne(ctx context.Context, paramID string) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, echokit.NewValidErr("paramID is missing")
	}
	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) findOne(ctx context.Context, id int64) (*mysqldb.Song, error) {
	books, err := b.Repo.Find(ctx, dbkit.Eq{mysqldb_repo.SongTable.ID: id})
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	}
	return books[0], nil
}

// Delete book
func (b *SongSvcImpl) Delete(ctx context.Context, paramID string) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return echokit.NewValidErr("paramID is missing")
	}
	_, err := b.Repo.Delete(ctx, dbkit.Eq{mysqldb_repo.SongTable.ID: id})
	return err
}

// Update book
func (b *SongSvcImpl) Update(ctx context.Context, paramID string, book *mysqldb.Song) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, echokit.NewValidErr("paramID is missing")
	}
	if err := validator.New().Struct(book); err != nil {
		return nil, echokit.NewValidErr(err.Error())
	}
	if _, err := b.findOne(ctx, id); err != nil {
		return nil, err
	}
	if err := b.update(ctx, id, book); err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) update(ctx context.Context, id int64, song *mysqldb.Song) error {
	affectedRow, err := b.Repo.Update(ctx, song, dbkit.Eq{mysqldb_repo.SongTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}

// Patch book
func (b *SongSvcImpl) Patch(ctx context.Context, paramID string, song *mysqldb.Song) (*mysqldb.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, echokit.NewValidErr("paramID is missing")
	}
	if _, err := b.findOne(ctx, id); err != nil {
		return nil, err
	}
	if err := b.patch(ctx, id, song); err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) patch(ctx context.Context, id int64, song *mysqldb.Song) error {
	affectedRow, err := b.Repo.Patch(ctx, song, dbkit.Eq{mysqldb_repo.SongTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}
