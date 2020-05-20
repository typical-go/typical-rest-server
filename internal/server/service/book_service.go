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
	// BookService contain logic for Book Controller
	// @mock
	BookService interface {
		FindOne(context.Context, string) (*repository.Book, error)
		Find(context.Context) ([]*repository.Book, error)
		Create(context.Context, *repository.Book) (int64, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *repository.Book) error
	}

	// BookServiceImpl is implementation of BookService
	BookServiceImpl struct {
		dig.In
		repository.BookRepo
	}
)

// NewBookService return new instance of BookService
// @constructor
func NewBookService(impl BookServiceImpl) BookService {
	return &impl
}

// Find books
func (b *BookServiceImpl) Find(ctx context.Context) ([]*repository.Book, error) {
	return b.BookRepo.Find(ctx)
}

// FindOne book
func (b *BookServiceImpl) FindOne(ctx context.Context, paramID string) (*repository.Book, error) {
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return nil, errvalid.Wrap(err)
	}

	books, err := b.BookRepo.Find(ctx, dbkit.Equal("id", id))
	if err != nil {
		return nil, err
	}

	if len(books) < 1 {
		return nil, sql.ErrNoRows
	}

	return books[0], nil
}

// Delete book
func (b *BookServiceImpl) Delete(ctx context.Context, paramID string) error {
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return errvalid.Wrap(err)
	}
	return b.BookRepo.Delete(ctx, id)
}

// Update book
func (b *BookServiceImpl) Update(ctx context.Context, paramID string, book *repository.Book) error {
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return errvalid.Wrap(err)
	}

	book.ID = id
	if err = validator.New().Struct(book); err != nil {
		return err
	}

	if _, err = b.BookRepo.Find(ctx, dbkit.Equal("id", id)); err != nil {
		return err
	}

	return b.BookRepo.Update(ctx, book)
}

// Create Book
func (b *BookServiceImpl) Create(ctx context.Context, book *repository.Book) (insertID int64, err error) {
	if err = validator.New().Struct(book); err != nil {
		return -1, err
	}
	return b.BookRepo.Create(ctx, book)
}
