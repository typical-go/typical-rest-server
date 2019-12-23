package repository

import (
	"context"
	"time"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

// Tag represented  tag entity
type Tag struct {
	ID         int64      `json:"id"`
	Type       string     `json:"type"`
	Attributes dbkit.JSON `json:"attributes"`
	Value      string     `json:"value"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TagRepo to handle tags entity
type TagRepo interface {
	Find(context.Context, int64) (*Tag, error)
	List(context.Context) ([]*Tag, error)
	Insert(context.Context, Tag) (lastInsertID int64, err error)
	Delete(context.Context, int64) error
	Update(context.Context, Tag) error
}

// NewTagRepo return new instance of TagRepo
func NewTagRepo(impl CachedTagRepoImpl) TagRepo {
	return &impl
}
