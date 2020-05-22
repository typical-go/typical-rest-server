package service

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"

	"github.com/typical-go/typical-rest-server/pkg/errvalid"
	"gopkg.in/go-playground/validator.v9"

	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"go.uber.org/dig"
)

type (
	// BookSvc contain logic for Book Controller
	// @mock
	BookSvc interface {
		FindOne(context.Context, string) (*repository.Book, error)
		Find(context.Context) ([]*repository.Book, error)
		Create(context.Context, *repository.Book) (int64, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *repository.Book) error
		Patch(context.Context, string, *repository.Book) error
	}

	// BookSvcImpl is implementation of BookSvc
	BookSvcImpl struct {
		dig.In
		repository.BookRepo
	}
)

// NewBookSvc return new instance of BookSvc
// @constructor
func NewBookSvc(impl BookSvcImpl) BookSvc {
	return &impl
}

// Find books
func (b *BookSvcImpl) Find(ctx context.Context) ([]*repository.Book, error) {
	return b.BookRepo.Find(ctx)
}

// FindOne book
func (b *BookSvcImpl) FindOne(ctx context.Context, paramID string) (*repository.Book, error) {
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return nil, errvalid.Wrap(err)
	}

	books, err := b.BookRepo.Find(ctx, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return nil, err
	}

	if len(books) < 1 {
		return nil, sql.ErrNoRows
	}

	return books[0], nil
}

// Delete book
func (b *BookSvcImpl) Delete(ctx context.Context, paramID string) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return errvalid.New("paramID is missing")
	}
	return b.BookRepo.Delete(
		ctx,
		dbkit.Equal(repository.BookCols.ID, id),
	)
}

// Update book
func (b *BookSvcImpl) Update(ctx context.Context, paramID string, book *repository.Book) (err error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return errvalid.New("paramID is missing")
	}

	if err = validator.New().Struct(book); err != nil {
		return err
	}

	if _, err = b.BookRepo.Find(
		ctx,
		dbkit.Equal(repository.BookCols.ID, id),
	); err != nil {
		return err
	}

	return b.BookRepo.Update(
		ctx,
		book,
		dbkit.Equal(repository.BookCols.ID, id),
	)
}

// Patch book
func (b *BookSvcImpl) Patch(ctx context.Context, paramID string, book *repository.Book) (err error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return errvalid.New("paramID is missing")
	}

	if err = validator.New().Struct(book); err != nil {
		return err
	}

	if _, err = b.BookRepo.Find(
		ctx,
		dbkit.Equal(repository.BookCols.ID, id),
	); err != nil {
		return err
	}

	return b.BookRepo.Patch(
		ctx,
		book,
		dbkit.Equal(repository.BookCols.ID, id),
	)
}

// Create Book
func (b *BookSvcImpl) Create(ctx context.Context, book *repository.Book) (insertID int64, err error) {
	if err = validator.New().Struct(book); err != nil {
		return -1, err
	}
	return b.BookRepo.Create(ctx, book)
}
