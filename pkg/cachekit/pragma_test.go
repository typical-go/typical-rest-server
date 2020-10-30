package cachekit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

func TestPragma_NoCache(t *testing.T) {
	testcases := []struct {
		testName string
		desc     string
		*cachekit.Pragma
		expected bool
	}{
		{
			desc:   "cache-control not available",
			Pragma: pragmaWithCacheControl(""),
		},
		{
			desc:     "lower case no-cache",
			Pragma:   pragmaWithCacheControl("no-cache"),
			expected: true,
		},
		{
			desc:     "upper case no-cache",
			Pragma:   pragmaWithCacheControl("NO-CACHE"),
			expected: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.NoCache, tt.desc)
		})
	}
}

func TestPragma_MaxAge(t *testing.T) {
	testcases := []struct {
		testName string
		desc     string
		*cachekit.Pragma
		expected time.Duration
	}{
		{
			desc:     "empty cache control",
			Pragma:   pragmaWithCacheControl(""),
			expected: cachekit.DefaultMaxAge,
		},
		{
			desc:     "empty cache control with new default max-age",
			Pragma:   pragmaWithCacheControl(""),
			expected: cachekit.DefaultMaxAge,
		},
		{
			desc:     "max-age in cache control",
			Pragma:   pragmaWithCacheControl("max-age=100"),
			expected: 100 * time.Second,
		},
		{
			desc:     "max-age is invalid type",
			Pragma:   pragmaWithCacheControl("max-age=invalid"),
			expected: cachekit.DefaultMaxAge,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.MaxAge, tt.desc)
		})
	}
}

func pragmaWithCacheControl(cacheControl string) *cachekit.Pragma {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(cachekit.HeaderCacheControl, cacheControl)
	return cachekit.CreatePragma(req)
}

func pragmaWithIfModifiedSince(ifModifiedSince string) *cachekit.Pragma {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(cachekit.HeaderIfModifiedSince, ifModifiedSince)
	return cachekit.CreatePragma(req)
}
