package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/restserver/repository"
	"go.uber.org/dig"
)

// BookService contain logic for Book Controller [mock]
type BookService interface {
	FindOne(context.Context, int64) (*repository.Book, error)
	Find(context.Context, ...dbkit.FindOption) ([]*repository.Book, error)
	Create(context.Context, *repository.Book) (*repository.Book, error)
	Delete(context.Context, int64) error
	Update(context.Context, *repository.Book) (*repository.Book, error)
}

// BookServiceImpl is implementation of BookService
type BookServiceImpl struct {
	dig.In
	repository.BookRepo

	Redis *redis.Client
}

// NewBookService return new instance of BookService [constructor]
func NewBookService(impl BookServiceImpl) BookService {
	return &impl
}

// FindOne book
func (r *BookServiceImpl) FindOne(ctx context.Context, id int64) (book *repository.Book, err error) {
	book = new(repository.Book)
	var (
		cacheStore = dbkit.NewCacheStore(r.Redis)
		key        = fmt.Sprintf("BOOK:FIND_ONE:%d", id)
	)
	cacheStore.Retrieve(ctx, key, book, func() (interface{}, error) {
		book, err = r.BookRepo.FindOne(ctx, id)
		return book, err
	})
	return
}

// Find books
func (r *BookServiceImpl) Find(ctx context.Context, opts ...dbkit.FindOption) (list []*repository.Book, err error) {
	var (
		cacheStore = dbkit.NewCacheStore(r.Redis)
		key        = fmt.Sprintf("BOOK:FIND:%+v", opts)
	)
	cacheStore.Retrieve(ctx, key, &list, func() (interface{}, error) {
		list, err := r.BookRepo.Find(ctx, opts...)
		return list, err
	})
	return
}
