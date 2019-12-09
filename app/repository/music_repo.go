package repository

import (
	"context"
	"time"
)

// Book represented database model
type Music struct {
	ID        int64     
	Title     string    
	Author    string    
	UpdatedAt time.Time 
	CreatedAt time.Time 
}

// BookRepo to get book data from databasesa
type Music Repo interface {
	Find(ctx context.Context, id int64) (*Music , error)
	List(ctx context.Context) ([]*Music , error)
	Insert(ctx context.Context, book Music ) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, book Music ) error
}
