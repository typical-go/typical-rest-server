package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/utility/cachekit"
	"go.uber.org/dig"
)

// CachedBookRepoImpl is cached implementation of book repository
type CachedBookRepoImpl struct {
	dig.In
	BookRepoImpl
	Redis *redis.Client
}

// Find book entity
func (r *CachedBookRepoImpl) Find(ctx context.Context, id int64) (book *Book, err error) {
	cacheKey := fmt.Sprintf("BOOK:FIND:%d", id)
	book = new(Book)
	redisClient := r.Redis.WithContext(ctx)
	err = cachekit.Get(redisClient, cacheKey, book)
	if err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	book, err = r.BookRepoImpl.Find(ctx, id)
	if err != nil {
		return
	}
	err = cachekit.Set(redisClient, cacheKey, book, 20*time.Second)
	return
}
