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

// CachedLocaleRepoImpl is cached implementation of locale repository
type CachedLocaleRepoImpl struct {
	dig.In
	LocaleRepoImpl
	Redis *redis.Client
}

// Find locale entity
func (r *CachedLocaleRepoImpl) Find(ctx context.Context, id int64) (e *Locale, err error) {
	cacheKey := fmt.Sprintf("LOCALES:FIND:%d", id)
	e = new(Locale)
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, e); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if e, err = r.LocaleRepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, e, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of locale entity
func (r *CachedLocaleRepoImpl) List(ctx context.Context) (list []*Locale, err error) {
	cacheKey := fmt.Sprintf("LOCALES:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.LocaleRepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
