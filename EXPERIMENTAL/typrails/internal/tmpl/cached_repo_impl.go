package tmpl

// CachedRepoImpl template
const CachedRepoImpl = `package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"go.uber.org/dig"
)

// Cached{{.Type}}RepoImpl is cached implementation of {{.Name}} repository
type Cached{{.Type}}RepoImpl struct {
	dig.In
	{{.Type}}RepoImpl
	Redis *redis.Client
}

// Find {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) FindOne(ctx context.Context, id int64) (e *{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:FIND:%d", id)
	e = new({{.Type}})
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, e); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if e, err = r.{{.Type}}RepoImpl.FindOne(ctx, id); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, e, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) Find(ctx context.Context) (list []*{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = dbkit.GetCache(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.{{.Type}}RepoImpl.Find(ctx); err != nil {
		return
	}
	if err2 := dbkit.SetCache(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
`
