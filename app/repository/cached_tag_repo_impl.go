package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"go.uber.org/dig"
)

// CachedTagRepoImpl is cached implementation of tag repository
type CachedTagRepoImpl struct {
	dig.In
	TagRepoImpl
	Redis *redis.Client
}

// Find tag entity
func (r *CachedTagRepoImpl) Find(ctx context.Context, id int64) (e *Tag, err error) {
	cacheKey := fmt.Sprintf("TAGS:FIND:%d", id)
	e = new(Tag)
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, e); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if e, err = r.TagRepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, e, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of tag entity
func (r *CachedTagRepoImpl) List(ctx context.Context) (list []*Tag, err error) {
	cacheKey := fmt.Sprintf("TAGS:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.TagRepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
