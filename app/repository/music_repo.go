package repository

import (
	"context"
	"time"
)

// Music represented  music entity
type Music struct {
	ID        int64     
	Title     string    
	Author    string    
	UpdatedAt time.Time 
	CreatedAt time.Time 
}

// MusicRepo to handle music  entity
type MusicRepo interface {
	Find(ctx context.Context, id int64) (*Music, error)
	List(ctx context.Context) ([]*Music, error)
	Insert(ctx context.Context, music Music) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, music Music) error
}
