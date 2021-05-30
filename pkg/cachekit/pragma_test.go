package cachekit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

func TestCreatePragma(t *testing.T) {
	testcases := []struct {
		testName string
		header   http.Header
		expected *cachekit.Pragma
	}{
		{
			header:   http.Header{},
			expected: &cachekit.Pragma{},
		},
		{
			header: newHeader(map[string]string{
				"If-Modified-Since": "Thu, 01 Dec 2020 16:00:00 GMT",
				"Cache-Control":     "no-cache,max-age=120",
			}),
			expected: &cachekit.Pragma{
				NoCache:         true,
				MaxAge:          120 * time.Second,
				IfModifiedSince: cachekit.ParseTime("Thu, 01 Dec 2020 16:00:00 GMT"),
			},
		},
		{
			header: newHeader(map[string]string{
				"Cache-Control": "no-cache,max-age=invalid",
			}),
			expected: &cachekit.Pragma{NoCache: true},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expected, cachekit.CreatePragma(tt.header))
		})
	}
}

func TestSetHeader(t *testing.T) {
	testcases := []struct {
		testName string
		pragma   *cachekit.Pragma
		expected http.Header
	}{
		{
			pragma: &cachekit.Pragma{
				NoCache:      true,
				MaxAge:       25 * time.Second,
				LastModified: cachekit.ParseTime("Thu, 01 Dec 2020 16:00:00 GMT"),
				Expires:      cachekit.ParseTime("Thu, 01 Dec 2020 16:00:25 GMT"),
			},
			expected: http.Header{
				"Cache-Control": []string{"no-cache"},
				"Expires":       []string{"Tue, 01 Dec 2020 16:00:25 GMT"},
				"Last-Modified": []string{"Tue, 01 Dec 2020 16:00:00 GMT"},
			},
		},
		{
			pragma:   &cachekit.Pragma{MaxAge: 25 * time.Second},
			expected: http.Header{"Cache-Control": []string{"max-age=25"}},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.pragma.Header())
		})
	}
}

func newHeader(m map[string]string) http.Header {
	header := make(http.Header)
	for k, v := range m {
		header.Add(k, v)
	}
	return header
}
