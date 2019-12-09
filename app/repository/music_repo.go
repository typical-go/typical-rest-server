package repository

import (
	"context"
	"time"
)

// Music represented database model
type Music struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// MusicRepo to get book data from databasesa
type MusicRepo interface {
	Find(ctx context.Context, id int64) (*Music, error)
	List(ctx context.Context) ([]*Music, error)
	Insert(ctx context.Context, music Music) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, music Music) error
}

// NewMusicRepo return new instance of MusicRepo
func NewMusicRepo(impl CachedMusicRepoImpl) MusicRepo {
	return &impl
}
