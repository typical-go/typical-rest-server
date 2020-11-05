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
	// Cached ...
	Cached struct {
		Bytes []byte
		Head  Head
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

const (
	suffixKeyTime = ":time"
	suffixKeyHead = ":head"
	suffixKeyBody = ":body"
)

// Middleware ...
func (s *Store) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()

		key := s.PrefixKey + req.URL.String()
		pragma := s.pragma(ctx, req.Header, key)

		if !pragma.LastModified.IsZero() {
			ifModifiedTime := pragma.IfModifiedSince
			if !ifModifiedTime.IsZero() && pragma.LastModified.Before(ifModifiedTime) {
				return echo.NewHTTPError(http.StatusNotModified)
			}

			if !pragma.NoCache {
				cached, err := s.getCached(ctx, key)
				if err == nil {
					resp := c.Response()
					addHeader(resp.Header(), pragma.Header())
					addHeader(resp.Header(), cached.Head.Header)
					resp.WriteHeader(cached.Head.StatusCode)
					resp.Write(cached.Bytes)
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

		addHeader(rec.HeaderMap, pragma.Header())
		copyResponseWriter(rec, ogResp)
		return nil
	}
}

func (s *Store) getCached(ctx context.Context, key string) (*Cached, error) {

	bytes, err := s.Client.Get(ctx, key+suffixKeyBody).Bytes()
	if err != nil {
		return nil, err
	}

	headerBytes, err := s.Client.Get(ctx, key+suffixKeyHead).Bytes()
	if err != nil {
		return nil, err
	}

	var head Head
	json.Unmarshal(headerBytes, &head)

	return &Cached{
		Bytes: bytes,
		Head:  head,
	}, nil
}

func (s *Store) store(ctx context.Context, key string, rec *httptest.ResponseRecorder, maxAge time.Duration) (time.Time, error) {
	lastModified := time.Now()
	headBytes, _ := json.Marshal(Head{
		StatusCode: rec.Code,
		Header:     rec.HeaderMap,
	})
	pipe := s.Client.TxPipeline()
	pipe.Set(ctx, key+suffixKeyTime, FormatTime(lastModified), maxAge)
	pipe.Set(ctx, key+suffixKeyBody, rec.Body.Bytes(), maxAge)
	pipe.Set(ctx, key+suffixKeyHead, string(headBytes), maxAge)

	_, err := pipe.Exec(ctx)
	return lastModified, err
}

func (s *Store) pragma(ctx context.Context, header http.Header, key string) *Pragma {
	pragma := CreatePragma(header)
	if pragma.MaxAge < 1 {
		pragma.MaxAge = s.DefaultMaxAge
	}

	lastModified := ParseTime(s.Client.Get(ctx, key+suffixKeyTime).Val())
	ttl, _ := s.Client.TTL(ctx, key+suffixKeyTime).Result()

	pragma.LastModified = lastModified
	pragma.Expires = time.Now().Add(ttl)
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

func copyResponseWriter(from *httptest.ResponseRecorder, to http.ResponseWriter) {
	for k := range from.HeaderMap {
		to.Header().Add(k, from.HeaderMap.Get(k))
	}
	to.WriteHeader(from.Code)
	to.Write(from.Body.Bytes()) // NOTE: commit the response
}

func addHeader(parent http.Header, header http.Header) {
	for k := range header {
		parent.Add(k, header.Get(k))
	}
}
