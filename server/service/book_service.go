package service

import (
	"context"
	"strconv"

	"github.com/typical-go/typical-rest-server/pkg/errvalid"
	"gopkg.in/go-playground/validator.v9"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/server/repository"
	"go.uber.org/dig"
)

// BookService contain logic for Book Controller
// @mock
type BookService interface {
	FindOne(context.Context, string) (*repository.Book, error)
	Find(context.Context, ...dbkit.FindOption) ([]*repository.Book, error)
	Create(context.Context, *repository.Book) (*repository.Book, error)
	Delete(context.Context, string) error
	Update(context.Context, string, *repository.Book) (*repository.Book, error)
}

// BookServiceImpl is implementation of BookService
type BookServiceImpl struct {
	dig.In

	repository.BookRepo
	Redis *redis.Client
}

// NewBookService return new instance of BookService
// @constructor
func NewBookService(impl BookServiceImpl) BookService {
	return &impl
}

// FindOne book
func (b *BookServiceImpl) FindOne(ctx context.Context, paramID string) (*repository.Book, error) {
	var (
		id  int64
		err error
	)
	if id, err = strconv.ParseInt(paramID, 10, 64); err != nil {
		return nil, errvalid.Wrap(err)
	}
	return b.BookRepo.FindOne(ctx, id)
}

// Delete book
func (b *BookServiceImpl) Delete(ctx context.Context, paramID string) error {
	var (
		id  int64
		err error
	)
	if id, err = strconv.ParseInt(paramID, 10, 64); err != nil {
		return errvalid.Wrap(err)
	}
	return b.BookRepo.Delete(ctx, id)
}

// Update book
func (b *BookServiceImpl) Update(ctx context.Context, paramID string, book *repository.Book) (*repository.Book, error) {
	var (
		id  int64
		err error
	)
	if id, err = strconv.ParseInt(paramID, 10, 64); err != nil {
		return nil, errvalid.Wrap(err)
	}
	if err = validator.New().Struct(book); err != nil {
		return nil, err
	}
	return b.BookRepo.Update(ctx, id, book)
}

// Create Book
func (b *BookServiceImpl) Create(ctx context.Context, form *repository.Book) (*repository.Book, error) {
	var (
		err error
	)
	if err = validator.New().Struct(form); err != nil {
		return nil, err
	}
	return nil, nil
}
