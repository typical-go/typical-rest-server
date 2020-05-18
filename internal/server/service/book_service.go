package service

import (
	"context"
	"strconv"

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
		Create(context.Context, *repository.Book) (*repository.Book, error)
		Delete(context.Context, string) error
		Update(context.Context, string, *repository.Book) (*repository.Book, error)
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
	return b.BookRepo.FindOne(ctx, id)
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
func (b *BookServiceImpl) Update(ctx context.Context, paramID string, book *repository.Book) (*repository.Book, error) {
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return nil, errvalid.Wrap(err)
	}
	if err = validator.New().Struct(book); err != nil {
		return nil, err
	}
	return b.BookRepo.Update(ctx, id, book)
}

// Create Book
func (b *BookServiceImpl) Create(ctx context.Context, book *repository.Book) (*repository.Book, error) {
	err := validator.New().Struct(book)
	if err != nil {
		return nil, err
	}
	return b.BookRepo.Create(ctx, book)
}
