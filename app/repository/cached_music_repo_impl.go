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

// CachedMusicRepoImpl is cached implementation of music repository
type CachedMusicRepoImpl struct {
	dig.In
	MusicRepoImpl
	Redis *redis.Client
}

// Find music entity
func (r *CachedMusicRepoImpl) Find(ctx context.Context, id int64) (e *Music, err error) {
	cacheKey := fmt.Sprintf("MUSICS:FIND:%d", id)
	e = new(Music)
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, e); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if e, err = r.MusicRepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, e, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of music entity
func (r *CachedMusicRepoImpl) List(ctx context.Context) (list []*Music, err error) {
	cacheKey := fmt.Sprintf("MUSICS:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.MusicRepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
