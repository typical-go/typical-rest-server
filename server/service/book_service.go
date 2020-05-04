package service

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/server/repository"
	"go.uber.org/dig"
)

type (
	// BookService contain logic for Book Controller
	// @mock
	BookService interface {
		FindOne(context.Context, int64) (*repository.Book, error)
		Find(context.Context, ...dbkit.FindOption) ([]*repository.Book, error)
		Create(context.Context, *repository.Book) (*repository.Book, error)
		Delete(context.Context, int64) error
		Update(context.Context, int64, *repository.Book) (*repository.Book, error)
	}

	// BookServiceImpl is implementation of BookService
	BookServiceImpl struct {
		dig.In

		repository.BookRepo
		Redis *redis.Client
	}
)

// NewBookService return new instance of BookService
// @constructor
func NewBookService(impl BookServiceImpl) BookService {
	return &impl
}
