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

// MusicRepo to handle musics entity [mock]
type MusicRepo interface {
	FindOne(context.Context, int64) (*Music, error)
	Find(context.Context) ([]*Music, error)
	Create(context.Context, Music) (lastInsertID int64, err error)
	Delete(context.Context, int64) error
	Update(context.Context, Music) error
}

// NewMusicRepo return new instance of MusicRepo [autowire]
func NewMusicRepo(impl CachedMusicRepoImpl) MusicRepo {
	return &impl
}
