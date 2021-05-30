package cachekit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type (
	// Store ...
	Store struct {
		*redis.Client
		DefaultMaxAge time.Duration
		PrefixKey     string
	}
	// CacheData cache data
	CacheData struct {
		LastModified string
		Body         []byte
		Head         Head
	}
	// Head ...
	Head struct {
		StatusCode int
		Header     http.Header
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

		key := s.PrefixKey + req.URL.String()
		pragma := CreatePragma(req.Header)
		if pragma.MaxAge < 1 {
			pragma.MaxAge = s.DefaultMaxAge
		}
		cacheData, err := s.getCacheData(ctx, key)
		if err == nil {
			ttl, _ := s.Client.TTL(ctx, key).Result()

			pragma.LastModified = ParseTime(cacheData.LastModified)
			pragma.Expires = time.Now().Add(ttl)

			if !pragma.LastModified.IsZero() {
				ifModifiedTime := pragma.IfModifiedSince
				if !ifModifiedTime.IsZero() && pragma.LastModified.Before(ifModifiedTime) {
					return echo.NewHTTPError(http.StatusNotModified)
				}

				if !pragma.NoCache {
					resp := c.Response()
					addHeader(resp.Header(), pragma.Header())
					addHeader(resp.Header(), cacheData.Head.Header)
					resp.WriteHeader(cacheData.Head.StatusCode)
					resp.Write(cacheData.Body)
					return nil
				}
			}

		}

		ogResp := c.Response()
		rec := httptest.NewRecorder()
		c.SetResponse(echo.NewResponse(rec, c.Echo()))

		if err := next(c); err != nil {
			c.SetResponse(ogResp)
			return err
		}

		lastModified, err := s.store(ctx, key, rec, pragma.MaxAge)
		if err != nil {
			c.SetResponse(ogResp)
			return err
		}

		pragma.LastModified = lastModified
		pragma.Expires = lastModified.Add(pragma.MaxAge)

		addHeader(rec.Header(), pragma.Header())
		copyResponseWriter(rec, ogResp)
		return nil
	}
}

func (s *Store) getCacheData(ctx context.Context, key string) (*CacheData, error) {
	bytes, err := s.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var data CacheData
	json.Unmarshal(bytes, &data)
	return &data, nil
}

func (s *Store) store(ctx context.Context, key string, rec *httptest.ResponseRecorder, maxAge time.Duration) (time.Time, error) {
	lastModified := time.Now()
	// FIXME: rec.HeaderMap is depracted
	bytes, _ := json.Marshal(CacheData{
		LastModified: FormatTime(lastModified),
		Body:         rec.Body.Bytes(),
		Head: Head{
			StatusCode: rec.Code,
			Header:     rec.Header(),
		},
	})

	return lastModified, s.Client.Set(ctx, key, bytes, maxAge).Err()
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

func copyResponseWriter(from *httptest.ResponseRecorder, to http.ResponseWriter) {
	for k := range from.Header() {
		to.Header().Add(k, from.Header().Get(k))
	}
	to.WriteHeader(from.Code)
	to.Write(from.Body.Bytes()) // NOTE: commit the response
}

func addHeader(parent http.Header, header http.Header) {
	for k := range header {
		parent.Add(k, header.Get(k))
	}
}
