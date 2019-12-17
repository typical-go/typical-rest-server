package repository

import (
	"context"
	"time"
)

// Music represented  music entity
type Music struct {
	ID        int64     `json:"id"`
	Artist    string    `json:"artist"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// MusicRepo to handle music  entity
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
