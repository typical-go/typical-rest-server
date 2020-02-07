package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"

	"github.com/go-redis/redis"

	"go.uber.org/dig"
)

// CachedBookRepoImpl is cached implementation of book repository
type CachedBookRepoImpl struct {
	dig.In
	BookRepoImpl
	Redis *redis.Client
}

// FindOne book
func (r *CachedBookRepoImpl) FindOne(ctx context.Context, id int64) (book *Book, err error) {
	cacheKey := fmt.Sprintf("BOOK:FIND:%d", id)
	book = new(Book)
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, book); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if book, err = r.BookRepoImpl.FindOne(ctx, id); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, book, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// Find books
func (r *CachedBookRepoImpl) Find(ctx context.Context, opts ...dbkit.FindOption) (list []*Book, err error) {
	cacheKey := fmt.Sprintf("BOOK:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.BookRepoImpl.Find(ctx, opts...); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
