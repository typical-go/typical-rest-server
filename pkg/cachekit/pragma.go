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
		IfModifiedSince time.Time
		LastModified    time.Time
		NoCache         bool
		MaxAge          time.Duration
		Expires         time.Time
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
		IfModifiedSince: ifModifiedSince,
		NoCache:         noCache,
		MaxAge:          maxAge,
	}
}

// SetHeader set header
func (c *Pragma) SetHeader(header http.Header) {

	if !c.Expires.IsZero() {
		header.Add(HeaderExpires, formatTime(c.Expires))
	}
	if !c.LastModified.IsZero() {
		header.Add(HeaderLastModified, formatTime(c.LastModified))
	}

	header.Add(HeaderCacheControl, c.String())
}

func (c *Pragma) String() string {
	var cc []string
	if c.NoCache {
		cc = append(cc, "no-cache")
	} else {
		cc = append(cc, fmt.Sprintf("max-age=%d", int(c.MaxAge.Seconds())))
	}
	return strings.Join(cc, " ")
}
