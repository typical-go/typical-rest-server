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
		RetrieveOne(context.Context, string) (*repository.Book, error)
		Retrieve(context.Context) ([]*repository.Book, error)
		Create(context.Context, *repository.Book) (*repository.Book, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *repository.Book) (*repository.Book, error)
		Patch(context.Context, string, *repository.Book) (*repository.Book, error)
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

// Create Book
func (b *BookSvcImpl) Create(ctx context.Context, book *repository.Book) (*repository.Book, error) {
	if err := validator.New().Struct(book); err != nil {
		return nil, errvalid.Wrap(err)
	}
	id, err := b.BookRepo.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	return b.retrieveOne(ctx, id)
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

	books, err := b.BookRepo.Retrieve(ctx, dbkit.Equal(repository.BookTable.ID, id))
	if err != nil {
		return nil, err
	} else if len(books) < 1 {
		return nil, sql.ErrNoRows
	}

	return books[0], nil
}

func (b *BookSvcImpl) retrieveOne(ctx context.Context, id int64) (*repository.Book, error) {
	books, err := b.BookRepo.Retrieve(ctx, dbkit.Equal(repository.BookTable.ID, id))
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

	_, err := b.BookRepo.Delete(ctx, dbkit.Equal(repository.BookTable.ID, id))
	return err
}

// Update book
func (b *BookSvcImpl) Update(ctx context.Context, paramID string, book *repository.Book) (*repository.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, errvalid.New("paramID is missing")
	}

	err := validator.New().Struct(book)
	if err != nil {
		return nil, errvalid.Wrap(err)
	}

	affectedRow, err := b.BookRepo.Update(ctx, book, dbkit.Equal(repository.BookTable.ID, id))
	if err != nil {
		return nil, err
	}
	if affectedRow < 1 {
		return nil, sql.ErrNoRows
	}

	return b.retrieveOne(ctx, id)
}

// Patch book
func (b *BookSvcImpl) Patch(ctx context.Context, paramID string, book *repository.Book) (*repository.Book, error) {
	id, _ := strconv.ParseInt(paramID, 10, 64)
	if id < 1 {
		return nil, errvalid.New("paramID is missing")
	}

	affectedRow, err := b.BookRepo.Patch(ctx, book, dbkit.Equal(repository.BookTable.ID, id))
	if err != nil {
		return nil, err
	}
	if affectedRow < 1 {
		return nil, sql.ErrNoRows
	}

	return b.retrieveOne(ctx, id)
}
