package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

// NewCacheStore return new instaence of cache store
// @ctor
func NewCacheStore(cfg *CacheCfg) *cachekit.Store {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logrus.Fatalf("redis: %s", err.Error())
	}

	return &cachekit.Store{
		Client:        client,
		DefaultMaxAge: cfg.DefaultMaxAge,
		PrefixKey:     cfg.PrefixKey,
	}
}
