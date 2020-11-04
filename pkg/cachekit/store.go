package cachekit

import (
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
)

var (
	gmt, _ = time.LoadLocation("GMT")
)

const (
	suffixKeyTime   = ":time"
	suffixKeyHeader = ":header"
)

// Middleware ...
func (s *Store) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		pragma := s.createPragma(req.Header)
		key := s.PrefixKey + req.URL.String()
		lastModified := ParseTime(s.Client.Get(ctx, key+suffixKeyTime).Val())

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
				headerBytes, err := s.Client.Get(ctx, key+suffixKeyHeader).Bytes()
				if err != nil {
					return err
				}

				pragma.LastModified = lastModified
				pragma.Expires = time.Now().Add(ttl)
				pragma.SetHeader(c.Response().Header())

				var header http.Header
				json.Unmarshal(headerBytes, &header)
				for k := range header {
					c.Response().Header().Add(k, header.Get(k))
				}
				c.Response().WriteHeader(http.StatusOK)
				c.Response().Write(data)
				return nil
			}
		}

		ogResp := c.Response()

		rec := httptest.NewRecorder()
		c.SetResponse(echo.NewResponse(rec, c.Echo()))

		if err := next(c); err != nil {
			c.SetResponse(ogResp)
			return err
		}

		maxAge := pragma.MaxAge
		lastModified = time.Now()

		pipe := s.Client.TxPipeline()
		pipe.Set(ctx, key, rec.Body.Bytes(), maxAge)
		pipe.Set(ctx, key+suffixKeyTime, FormatTime(lastModified), maxAge)

		headerBytes, _ := json.Marshal(rec.HeaderMap)
		pipe.Set(ctx, key+suffixKeyHeader, string(headerBytes), maxAge)
		if _, err := pipe.Exec(ctx); err != nil {
			c.SetResponse(ogResp)
			return err
		}

		pragma.LastModified = lastModified
		pragma.Expires = lastModified.Add(maxAge)
		pragma.SetHeader(c.Response().Header())

		copyResponseWriter(rec, ogResp)
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

func copyResponseWriter(from *httptest.ResponseRecorder, to http.ResponseWriter) {
	for k := range from.HeaderMap {
		to.Header().Add(k, from.HeaderMap.Get(k))
	}
	to.WriteHeader(from.Code)
	to.Write(from.Body.Bytes()) // NOTE: commit the response
}
