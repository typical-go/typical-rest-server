package cachekit

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultMaxAge is default max-age value
	DefaultMaxAge = 30 * time.Second

	// HeaderCacheControl as in https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	HeaderCacheControl = "Cache-Control"

	// HeaderExpires as in https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Expires
	HeaderExpires = "Expires"

	// HeaderLastModified as in https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified
	HeaderLastModified = "Last-Modified"

	// HeaderIfModifiedSince as in https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since
	HeaderIfModifiedSince = "If-Modified-Since"
)

type (
	// Pragma handle pragmatic information/directives for caching
	Pragma struct {
		ifModifiedSince time.Time
		lastModified    time.Time
		noCache         bool
		maxAge          time.Duration
		expires         time.Time
	}
)

// CreatePragma to create new instance of CacheControl from request
func CreatePragma(req *http.Request) *Pragma {
	var noCache bool
	var ifModifiedSince time.Time
	maxAge := DefaultMaxAge

	if raw := req.Header.Get(HeaderIfModifiedSince); raw != "" {
		ifModifiedSince, _ = time.Parse(time.RFC1123, raw)
	}

	if raw := req.Header.Get(HeaderCacheControl); raw != "" {
		for _, s := range strings.Split(raw, ",") {
			s = strings.ToLower(strings.TrimSpace(s))
			if s == "no-cache" {
				noCache = true
			}

			maxAgeField := "max-age="
			if strings.HasPrefix(s, maxAgeField) {
				maxAgeRaw, err := strconv.Atoi(s[len(maxAgeField):])
				if err != nil {
					break
				}
				maxAge = time.Duration(maxAgeRaw) * time.Second
			}
		}
	}
	return &Pragma{
		ifModifiedSince: ifModifiedSince,
		noCache:         noCache,
		maxAge:          maxAge,
	}
}

// SetExpiresByTTL to set expires to current time to TTL
func (c *Pragma) SetExpiresByTTL(ttl time.Duration) {
	c.expires = time.Now().Add(ttl)
}

// SetLastModified to set last-modified
func (c *Pragma) SetLastModified(lastModified time.Time) {
	c.lastModified = lastModified
}

// NoCache return true if no cache is set
func (c *Pragma) NoCache() bool {
	return c.noCache
}

// MaxAge return max-age cache (in seconds)
func (c *Pragma) MaxAge() time.Duration {
	return c.maxAge
}

// IfModifiedSince return if-modified-since value
func (c *Pragma) IfModifiedSince() time.Time {
	return c.ifModifiedSince
}

// ResponseHeaders return map that contain response header
func (c *Pragma) ResponseHeaders() map[string]string {
	m := make(map[string]string)
	if !c.expires.IsZero() {
		m[HeaderExpires] = GMT(c.expires).Format(time.RFC1123)
	}
	if !c.lastModified.IsZero() {
		m[HeaderLastModified] = GMT(c.lastModified).Format(time.RFC1123)
	}
	if cacheControls := c.respCacheControls(); len(cacheControls) > 0 {
		m[HeaderCacheControl] = strings.Join(cacheControls, " ")
	}
	return m
}

func (c *Pragma) respCacheControls() (cc []string) {
	if c.NoCache() {
		cc = append(cc, "no-cache")
	} else {
		cc = append(cc, fmt.Sprintf("max-age=%d", int(c.MaxAge().Seconds())))
	}
	return
}
