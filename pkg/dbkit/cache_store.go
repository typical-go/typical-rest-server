package dbkit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/copier"
)

// CacheStore responsible to cache data
type CacheStore struct {
	*redis.Client
	expiration time.Duration
}

// NewCacheStore return new instance of CacheStore
func NewCacheStore(client *redis.Client) *CacheStore {
	return &CacheStore{
		Client:     client,
		expiration: 100 * time.Millisecond,
	}
}

// WithExpiration return cache store with new expiration
func (c *CacheStore) WithExpiration(expiration time.Duration) *CacheStore {
	c.expiration = expiration
	return c
}

// Retrieve cache data
func (c *CacheStore) Retrieve(ctx context.Context, key string, target interface{}, refresh func() (interface{}, error)) (err error) {
	var data []byte
	if data, err = c.WithContext(ctx).Get(key).Bytes(); err != nil {
		// NOTE: cache not available
		var v interface{}
		if v, err = refresh(); err != nil {
			return
		}
		if data, err = c.marshal(v); err != nil {
			return
		}
		if err = c.WithContext(ctx).Set(key, data, c.expiration).Err(); err != nil {
			return
		}

		copier.Copy(target, v)
		return
	}
	return c.unmarshal(data, target)
}

func (c *CacheStore) marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *CacheStore) unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
