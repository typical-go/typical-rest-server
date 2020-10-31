package cachekit

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

type (
	// Store ...
	Store struct {
		Client        *redis.Client
		DefaultMaxAge time.Duration
		Prefix        string
	}
)

var (
	gmt, _ = time.LoadLocation("GMT")
)

// Middleware ...
func (s *Store) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		pragma := s.createPragma(req.Header)
		key := s.Prefix + req.URL.String()
		lastModified := ParseTime(s.Client.Get(ctx, key+":time").Val())

		if !lastModified.IsZero() {
			ifModifiedTime := pragma.IfModifiedSince
			if !ifModifiedTime.IsZero() && lastModified.Before(ifModifiedTime) {
				return echo.NewHTTPError(http.StatusNotModified)
			}

			if !pragma.NoCache {
				ttl, err := s.Client.TTL(ctx, key).Result()
				if err != nil {
					return err
				}

				data, err := s.Client.Get(ctx, key).Bytes()
				if err != nil {
					return err
				}
				contentType, err := s.Client.Get(ctx, key+":type").Bytes()
				if err != nil {
					return err
				}

				pragma.LastModified = lastModified
				pragma.Expires = time.Now().Add(ttl)
				pragma.SetHeader(c.Response().Header())
				c.Response().Header().Add("Content-Type", string(contentType))
				c.Response().WriteHeader(http.StatusOK)
				c.Response().Write(data)
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

		pipe := s.Client.TxPipeline()
		pipe.Set(ctx, key, rw.Bytes, maxAge).Err()
		pipe.Set(ctx, key+":time", FormatTime(lastModified), maxAge).Err()
		pipe.Set(ctx, key+":type", rw.Header().Get("Content-Type"), maxAge).Err()

		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}

		pragma.LastModified = lastModified
		pragma.Expires = lastModified.Add(maxAge)
		pragma.SetHeader(c.Response().Header())

		rw.CopyTo(c.Response())

		return nil
	}
}

func (s *Store) createPragma(header http.Header) *Pragma {
	pragma := CreatePragma(header)
	if pragma.MaxAge < 1 {
		pragma.MaxAge = s.DefaultMaxAge
	}
	return pragma
}

// FormatTime format time
func FormatTime(t time.Time) string {
	return t.In(gmt).Format(time.RFC1123)
}

// ParseTime parse time
func ParseTime(raw string) time.Time {
	t, _ := time.Parse(time.RFC1123, raw)
	return t
}
