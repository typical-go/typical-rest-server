package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

type (
	// CacheCfg cache onfiguration
	// @envconfig (prefix:"CACHE")
	CacheCfg struct {
		DefaultMaxAge time.Duration `envconfig:"DEFAULT_MAX_AGE" default:"30s"`
		PrefixKey     string        `envconfig:"PREFIX_KEY" default:"cache_"`
		RedisHost     string        `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
		RedisPort     string        `envconfig:"REDIS_PORT" required:"true" default:"6379"`
		RedisPass     string        `envconfig:"REDIS_PASS" default:"redispass"`
	}
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
