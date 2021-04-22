package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/typical-go/typical-rest-server/internal/app/entity"
	"github.com/typical-go/typical-rest-server/internal/generated/dbrepo"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// SongSvc contain logic for Song Controller
	// @mock
	SongSvc interface {
		FindOne(context.Context, string) (*entity.Song, error)
		Find(context.Context, *FindSongReq) (*FindSongResp, error)
		Create(context.Context, *entity.Song) (*entity.Song, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *entity.Song) (*entity.Song, error)
		Patch(context.Context, string, *entity.Song) (*entity.Song, error)
	}
	// SongSvcImpl is implementation of SongSvc
	SongSvcImpl struct {
		dig.In
		Repo dbrepo.SongRepo
	}
	// FindSongReq find request
	FindSongReq struct {
		Limit  uint64 `query:"limit"`
		Offset uint64 `query:"offset"`
		Sort   string `query:"sort"`
	}
	// FindSongResp find song response
	FindSongResp struct {
		Songs      []*entity.Song
		TotalCount string
	}
)

// NewSongSvc return new instance of SongSvc
// @ctor
func NewSongSvc(impl SongSvcImpl) SongSvc {
	return &impl
}

// Create Song
func (b *SongSvcImpl) Create(ctx context.Context, book *entity.Song) (*entity.Song, error) {
	if err := validator.New().Struct(book); err != nil {
		return nil, echokit.NewValidErr(err.Error())
	}
	id, err := b.Repo.Insert(ctx, book)
	if err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

// Find books
func (b *SongSvcImpl) Find(ctx context.Context, req *FindSongReq) (*FindSongResp, error) {
	var opts []sqkit.SelectOption
	opts = append(opts, &sqkit.OffsetPagination{Offset: req.Offset, Limit: req.Limit})
	if req.Sort != "" {
		opts = append(opts, sqkit.Sorts(strings.Split(req.Sort, ",")))
	}
	totalCount, err := b.Repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	songs, err := b.Repo.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &FindSongResp{
		TotalCount: fmt.Sprintf("%d", totalCount),
		Songs:      songs,
	}, nil
}

// FindOne book
func (b *SongSvcImpl) FindOne(ctx context.Context, paramID string) (*entity.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)

	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) findOne(ctx context.Context, id int64) (*entity.Song, error) {
	books, err := b.Repo.Find(ctx, sqkit.Eq{dbrepo.SongTable.ID: id})
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
		return nil, echo.ErrNotFound
	}
	return books[0], nil
}

// Delete book
func (b *SongSvcImpl) Delete(ctx context.Context, paramID string) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	_, err := b.Repo.Delete(ctx, sqkit.Eq{dbrepo.SongTable.ID: id})
	return err
}

// Update book
func (b *SongSvcImpl) Update(ctx context.Context, paramID string, book *entity.Song) (*entity.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)

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

func (b *SongSvcImpl) update(ctx context.Context, id int64, song *entity.Song) error {
	affectedRow, err := b.Repo.Update(ctx, song, sqkit.Eq{dbrepo.SongTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}

// Patch book
func (b *SongSvcImpl) Patch(ctx context.Context, paramID string, song *entity.Song) (*entity.Song, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)

	if _, err := b.findOne(ctx, id); err != nil {
		return nil, err
	}
	if err := b.patch(ctx, id, song); err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *SongSvcImpl) patch(ctx context.Context, id int64, song *entity.Song) error {
	affectedRow, err := b.Repo.Patch(ctx, song, sqkit.Eq{dbrepo.SongTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}
