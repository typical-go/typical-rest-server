package cachekit

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	gmt, _ = time.LoadLocation("GMT")
)

// NotModifiedError return true if error is not modified error
func NotModifiedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not modified")
}

// SetHeader to set cache-related http header
func SetHeader(resp *echo.Response, pragma *Pragma) {
	header := resp.Header()
	for key, value := range pragma.ResponseHeaders() {
		header.Set(key, value)
	}
}

// GMT to convert time to GMT
func GMT(t time.Time) time.Time {
	return t.In(gmt)
}
