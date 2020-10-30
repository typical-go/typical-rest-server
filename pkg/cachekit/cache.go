package cachekit

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/copier"
)

var (
	// ErrNotModified happen when conditional request apply
	ErrNotModified = errors.New("cache: not modified")
)

type (
	// Cache data
	Cache struct {
		Client    *redis.Client
		Key       string
		RefreshFn RefreshFn
	}
	// RefreshFn is function that retrieve refresh data
	RefreshFn func() (interface{}, error)
)

// Execute cache to retreive data and save to target variable
func (c *Cache) Execute(target interface{}, pragma *Pragma) error {
	lastModified, err := c.getModifiedTime()
	if err != nil {
		return err
	}
	pragma.SetLastModified(lastModified)

	if !lastModified.IsZero() {
		ifModifiedTime := pragma.IfModifiedSince()
		if !ifModifiedTime.IsZero() && lastModified.Before(ifModifiedTime) {
			return ErrNotModified
		}

		if !pragma.NoCache() {
			ttl, err := c.getCached(target)
			if err != nil {
				return err
			}
			pragma.SetExpiresByTTL(ttl)
			return nil
		}
	}

	v, err := c.RefreshFn()
	if err != nil {
		return err
	}

	ttl := pragma.MaxAge()
	modifiedTime, err := c.store(v, ttl)
	if err != nil {
		return err
	}

	pragma.SetLastModified(modifiedTime)
	pragma.SetExpiresByTTL(ttl)

	return copier.Copy(target, v)

}

func (c *Cache) store(v interface{}, ttl time.Duration) (time.Time, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return time.Time{}, err
	}
	err = c.Client.Set(c.Key, data, ttl).Err()
	if err != nil {
		return time.Time{}, err
	}

	modifiedTime := time.Now()
	err = c.Client.Set(c.modifiedTimeKey(), GMT(modifiedTime).Format(time.RFC1123), ttl).Err()
	if err != nil {
		return time.Time{}, err
	}
	return modifiedTime, nil
}

func (c *Cache) getCached(target interface{}) (time.Duration, error) {
	ttl, err := c.Client.TTL(c.Key).Result()
	if err != nil {
		return 0, err
	}
	data, err := c.Client.Get(c.Key).Bytes()
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, target); err != nil {
		return 0, err
	}
	return ttl, err
}

func (c *Cache) getModifiedTime() (time.Time, error) {
	raw := c.Client.Get(c.modifiedTimeKey()).Val()
	if raw == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC1123, raw)
}

func (c *Cache) modifiedTimeKey() string {
	return c.Key + ":time"
}
