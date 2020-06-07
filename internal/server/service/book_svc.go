package service

import (
	"context"
	"database/sql"
	"errors"
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
		RetrieveOne(context.Context, string) (*repository.Book, error)
		Retrieve(context.Context) ([]*repository.Book, error)
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
// @ctor
func NewBookSvc(impl BookSvcImpl) BookSvc {
	return &impl
}

// Retrieve books
func (b *BookSvcImpl) Retrieve(ctx context.Context) ([]*repository.Book, error) {
	return b.BookRepo.Retrieve(ctx)
}

// RetrieveOne book
func (b *BookSvcImpl) RetrieveOne(ctx context.Context, paramID string) (*repository.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, errvalid.New("paramID is missing")
	}

	books, err := b.BookRepo.Retrieve(ctx, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
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

	affectedRow, err := b.BookRepo.Delete(ctx, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return err
	} else if affectedRow < 1 {
		return errors.New("No deleted row")
	}

	return nil
}

// Update book
func (b *BookSvcImpl) Update(ctx context.Context, paramID string, book *repository.Book) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return errvalid.New("paramID is missing")
	}

	err := validator.New().Struct(book)
	if err != nil {
		return err
	}

	_, err = b.BookRepo.Retrieve(ctx, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return err
	}

	affectedRow, err := b.BookRepo.Update(ctx, book, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return err
	} else if affectedRow < 1 {
		return errors.New("No updated row")
	}

	return nil
}

// Patch book
func (b *BookSvcImpl) Patch(ctx context.Context, paramID string, book *repository.Book) error {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return errvalid.New("paramID is missing")
	}

	_, err := b.BookRepo.Retrieve(ctx, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return err
	}

	affectedRow, err := b.BookRepo.Patch(ctx, book, dbkit.Equal(repository.BookCols.ID, id))
	if err != nil {
		return err
	} else if affectedRow < 1 {
		return errors.New("No patched row")
	}

	return nil
}

// Create Book
func (b *BookSvcImpl) Create(ctx context.Context, book *repository.Book) (int64, error) {
	err := validator.New().Struct(book)
	if err != nil {
		return -1, err
	}
	return b.BookRepo.Create(ctx, book)
}
