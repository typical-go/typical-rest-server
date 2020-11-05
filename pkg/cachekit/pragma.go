package cachekit

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
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
func CreatePragma(header http.Header) *Pragma {
	var noCache bool
	var ifModifiedSince time.Time
	var maxAge time.Duration

	ifModifiedSince = ParseTime(header.Get(HeaderIfModifiedSince))

	if raw := header.Get(HeaderCacheControl); raw != "" {
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

// Header set header
func (c *Pragma) Header() http.Header {
	header := make(http.Header)
	if !c.Expires.IsZero() {
		header.Add(HeaderExpires, FormatTime(c.Expires))
	}
	if !c.LastModified.IsZero() {
		header.Add(HeaderLastModified, FormatTime(c.LastModified))
	}
	header.Add(HeaderCacheControl, c.String())
	return header
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
