package infra

import (
	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

// NewCacheStore return new instaence of cache store
// @ctor
func NewCacheStore(client *redis.Client) *cachekit.Store {
	return &cachekit.Store{Client: client}
}
