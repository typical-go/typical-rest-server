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
)

type (
	// BookSvc contain logic for Book Controller
	// @mock
	BookSvc interface {
		FindOne(context.Context, string) (*entity.Book, error)
		Find(context.Context, *FindBookReq) (*FindBookResp, error)
		Create(context.Context, *entity.Book) (*entity.Book, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *entity.Book) (*entity.Book, error)
		Patch(context.Context, string, *entity.Book) (*entity.Book, error)
	}
	// BookSvcImpl is implementation of BookSvc
	BookSvcImpl struct {
		dig.In
		Repo dbrepo.BookRepo
	}
	// FindBookReq find request
	FindBookReq struct {
		Limit  uint64 `query:"limit"`
		Offset uint64 `query:"offset"`
		Sort   string `query:"sort"`
	}
	// FindBookResp find book resp
	FindBookResp struct {
		Books      []*entity.Book
		TotalCount string
	}
)

// NewBookSvc return new instance of BookSvc
// @ctor
func NewBookSvc(impl BookSvcImpl) BookSvc {
	return &impl
}

// Create Book
func (b *BookSvcImpl) Create(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	if errMsg := b.validateBook(book); errMsg != "" {
		return nil, echokit.NewValidErr(errMsg)
	}
	id, err := b.Repo.Insert(ctx, book)
	if err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *BookSvcImpl) validateBook(book *entity.Book) string {
	if book.Title == "" {
		return "title must be filled"
	}
	if book.Author == "" {
		return "author must be filled"
	}
	return ""
}

// Find books
func (b *BookSvcImpl) Find(ctx context.Context, req *FindBookReq) (*FindBookResp, error) {
	var opts []sqkit.SelectOption
	opts = append(opts, &sqkit.OffsetPagination{Offset: req.Offset, Limit: req.Limit})
	if req.Sort != "" {
		opts = append(opts, sqkit.Sorts(strings.Split(req.Sort, ",")))
	}
	totalCount, err := b.Repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	books, err := b.Repo.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &FindBookResp{
		Books:      books,
		TotalCount: fmt.Sprintf("%d", totalCount),
	}, nil
}

// FindOne book
func (b *BookSvcImpl) FindOne(ctx context.Context, paramID string) (*entity.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	return b.findOne(ctx, id)
}

func (b *BookSvcImpl) findOne(ctx context.Context, id int64) (*entity.Book, error) {
	books, err := b.Repo.Find(ctx, sqkit.Eq{dbrepo.BookTable.ID: id})
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
		return nil, echo.ErrNotFound
	}
	return books[0], nil
}

// Delete book
func (b *BookSvcImpl) Delete(ctx context.Context, paramID string) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	_, err := b.Repo.Delete(ctx, sqkit.Eq{dbrepo.BookTable.ID: id})
	return err
}

// Update book
func (b *BookSvcImpl) Update(ctx context.Context, paramID string, book *entity.Book) (*entity.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if errMsg := b.validateBook(book); errMsg != "" {
		return nil, echokit.NewValidErr(errMsg)
	}
	if _, err := b.findOne(ctx, id); err != nil {
		return nil, err
	}
	if err := b.update(ctx, id, book); err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *BookSvcImpl) update(ctx context.Context, id int64, book *entity.Book) error {
	affectedRow, err := b.Repo.Update(ctx, book, sqkit.Eq{dbrepo.BookTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}

// Patch book
func (b *BookSvcImpl) Patch(ctx context.Context, paramID string, book *entity.Book) (*entity.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if _, err := b.findOne(ctx, id); err != nil {
		return nil, err
	}
	if err := b.patch(ctx, id, book); err != nil {
		return nil, err
	}
	return b.findOne(ctx, id)
}

func (b *BookSvcImpl) patch(ctx context.Context, id int64, book *entity.Book) error {
	affectedRow, err := b.Repo.Patch(ctx, book, sqkit.Eq{dbrepo.BookTable.ID: id})
	if err != nil {
		return err
	}
	if affectedRow < 1 {
		return errors.New("no affected row")
	}
	return nil
}
