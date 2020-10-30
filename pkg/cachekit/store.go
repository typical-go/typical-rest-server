package cachekit

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

type (
	// Store ...
	Store struct {
		Client *redis.Client
	}
)

var (
	gmt, _ = time.LoadLocation("GMT")
)

// Middleware ...
func (s *Store) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		pragma := CreatePragma(req)
		key := req.URL.String()

		lastModified, err := s.getLastModified(key)
		if err != nil {
			return err
		}

		if !lastModified.IsZero() {
			ifModifiedTime := pragma.IfModifiedSince
			if !ifModifiedTime.IsZero() && lastModified.Before(ifModifiedTime) {
				return echo.NewHTTPError(http.StatusNotModified)
			}

			if !pragma.NoCache {
				data, ttl, err := s.get(key)
				if err != nil {
					return err
				}
				pragma.LastModified = lastModified
				pragma.Expires = time.Now().Add(ttl)
				pragma.SetHeader(c.Response().Header())
				c.Response().Write(data)
				c.Response().WriteHeader(http.StatusOK)
				return nil
			}
		}

		rw := echokit.NewResponseWriter()
		c1 := c.Echo().NewContext(req, rw)

		if err := next(c1); err != nil {
			c.Error(err)
			return err
		}

		maxAge := pragma.MaxAge
		lastModified = time.Now()

		if err := s.set(key, rw.Bytes, maxAge); err != nil {
			return err
		}
		if err := s.setLastModifid(key, lastModified, maxAge); err != nil {
			return err
		}
		pragma.LastModified = lastModified
		pragma.Expires = lastModified.Add(maxAge)
		pragma.SetHeader(c.Response().Header())

		rw.CopyTo(c.Response())

		return nil
	}
}

func (s *Store) getLastModified(key string) (time.Time, error) {
	raw := s.Client.Get(key + ":time").Val()
	if raw == "" {
		return time.Time{}, nil
	}
	return parseTime(raw)
}

func (s *Store) setLastModifid(key string, lastModified time.Time, expr time.Duration) error {
	return s.Client.Set(key+":time", formatTime(lastModified), expr).Err()
}

func (s *Store) set(key string, b []byte, expr time.Duration) error {
	return s.Client.Set(key, b, expr).Err()
}

func (s *Store) get(key string) ([]byte, time.Duration, error) {
	ttl, err := s.Client.TTL(key).Result()
	if err != nil {
		return nil, 0, err
	}
	data, err := s.Client.Get(key).Bytes()
	if err != nil {
		return nil, 0, err
	}
	return data, ttl, nil
}

func (s *Store) getCached(key string) ([]byte, error) {
	return s.Client.Get(key).Bytes()
}

func formatTime(t time.Time) string {
	return t.In(gmt).Format(time.RFC1123)
}

func parseTime(raw string) (time.Time, error) {
	return time.Parse(time.RFC1123, raw)
}
