package repository

import (
	"context"
	"time"
)

// Book represented database model
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// BookRepo to get book data from databasesa
type BookRepo interface {
	Find(ctx context.Context, id int64) (*Book, error)
	List(ctx context.Context) ([]*Book, error)
	Insert(ctx context.Context, book Book) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, book Book) error
}

// NewBookRepo return new instance of BookRepo
func NewBookRepo(impl CachedBookRepoImpl) BookRepo {
	return &impl
}
