package tmpl

// CachedRepoImpl template
const CachedRepoImpl = `package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/utility/cachekit"
	"go.uber.org/dig"
)

// Cached{{.Type}}RepoImpl is cached implementation of {{.Name}} repository
type Cached{{.Type}}RepoImpl struct {
	dig.In
	{{.Type}}RepoImpl
	Redis *redis.Client
}

// Find {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) Find(ctx context.Context, id int64) ({{.Name}} *{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:FIND:%d", id)
	{{.Name}} = new({{.Type}})
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, {{.Name}}); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if {{.Name}}, err = r.{{.Type}}RepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, {{.Name}}, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) List(ctx context.Context) (list []*{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.{{.Type}}RepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
`
