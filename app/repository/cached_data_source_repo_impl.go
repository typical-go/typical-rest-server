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

// CachedDataSourceRepoImpl is cached implementation of data_source repository
type CachedDataSourceRepoImpl struct {
	dig.In
	DataSourceRepoImpl
	Redis *redis.Client
}

// Find data_source entity
func (r *CachedDataSourceRepoImpl) Find(ctx context.Context, id int64) (e *DataSource, err error) {
	cacheKey := fmt.Sprintf("DATA_SOURCES:FIND:%d", id)
	e = new(DataSource)
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, e); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if e, err = r.DataSourceRepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, e, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of data_source entity
func (r *CachedDataSourceRepoImpl) List(ctx context.Context) (list []*DataSource, err error) {
	cacheKey := fmt.Sprintf("DATA_SOURCES:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.DataSourceRepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
