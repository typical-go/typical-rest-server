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
	if err = cachekit.Get(redisClient, cacheKey, book); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if book, err = r.BookRepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, book, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of book entity
func (r *CachedBookRepoImpl) List(ctx context.Context) (list []*Book, err error) {
	cacheKey := fmt.Sprintf("BOOK:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.BookRepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
