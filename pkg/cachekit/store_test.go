package cachekit_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/typical-go/typical-rest-server/pkg/echokit"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

func TestStore_Middleware(t *testing.T) {
	defer monkey.Patch(time.Now, func() time.Time {
		return time.Date(2020, time.December, 16, 0, 0, 0, 0, time.UTC)
	}).Unpatch()

	testcases := []struct {
		testName      string
		next          echo.HandlerFunc
		defaultMaxAge time.Duration
		prefix        string
		header        map[string]string
		beforeFn      func(*miniredis.Miniredis)
		assertFn      func(*testing.T, *miniredis.Miniredis)
		expected      *echokit.ResponseWriter
		expectedErr   string
	}{
		{
			testName: "no cache available",
			next: func(ec echo.Context) error {
				return ec.JSON(200, "some-response")
			},
			defaultMaxAge: 30 * time.Second,
			prefix:        "cache_",
			expected: &echokit.ResponseWriter{
				StatusCode: 200,
				Bytes:      []byte("\"some-response\"\n"),
				RespHeader: http.Header{
					"Cache-Control": []string{"max-age=30"},
					"Content-Type":  []string{"application/json; charset=UTF-8"},
					"Expires":       []string{"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": []string{"Wed, 16 Dec 2020 00:00:00 GMT"},
				},
			},
			assertFn: func(t *testing.T, r *miniredis.Miniredis) {
				data, _ := r.Get("cache_/")
				lastModified, _ := r.Get("cache_/:time")
				contentType, _ := r.Get("cache_/:type")

				require.Equal(t, "\"some-response\"\n", data)
				require.Equal(t, "Wed, 16 Dec 2020 00:00:00 GMT", lastModified)
				require.Equal(t, "application/json; charset=UTF-8", contentType)

				require.Equal(t, 30*time.Second, r.TTL("cache_/"))
				require.Equal(t, 30*time.Second, r.TTL("cache_/:time"))
				require.Equal(t, 30*time.Second, r.TTL("cache_/:type"))
			},
		},
		{
			testName: "no cache available, got error",
			next: func(ec echo.Context) error {
				return errors.New("some-error")
			},
			defaultMaxAge: 30 * time.Second,
			expectedErr:   "some-error",
		},
		{
			testName:      "cache available",
			defaultMaxAge: 30 * time.Second,
			prefix:        "cache_",
			expected: &echokit.ResponseWriter{
				StatusCode: 200,
				Bytes:      []byte("\"some-response\"\n"),
				RespHeader: http.Header{
					"Cache-Control": []string{"max-age=30"},
					"Content-Type":  []string{"some-type"},
					"Expires":       []string{"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": []string{"Wed, 16 Dec 2020 00:00:00 GMT"},
				},
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "\"some-response\"\n")
				r.SetTTL("cache_/", 30*time.Second)
				r.Set("cache_/:time", "Wed, 16 Dec 2020 00:00:00 GMT")
				r.SetTTL("cache_/:time", 30*time.Second)
				r.Set("cache_/:type", "some-type")
				r.SetTTL("cache_/:type", 30*time.Second)
			},
		},
		{
			testName:      "if cache not modified since",
			defaultMaxAge: 30 * time.Second,
			prefix:        "cache_",
			header: map[string]string{
				"If-Modified-Since": "Wed, 16 Dec 2020 00:00:05 GMT",
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "\"some-response\"\n")
				r.SetTTL("cache_/", 30*time.Second)
				r.Set("cache_/:time", "Wed, 16 Dec 2020 00:00:00 GMT")
				r.SetTTL("cache_/:time", 30*time.Second)
				r.Set("cache_/:type", "some-type")
				r.SetTTL("cache_/:type", 30*time.Second)
			},
			expectedErr: "code=304, message=Not Modified",
		},
		{
			testName:      "if cache modified since",
			defaultMaxAge: 30 * time.Second,
			prefix:        "cache_",
			next: func(ec echo.Context) error {
				return ec.JSON(200, "some-response")
			},
			header: map[string]string{
				"If-Modified-Since": "Wed, 15 Dec 2020 23:59:00 GMT",
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "\"some-response\"\n")
				r.SetTTL("cache_/", 30*time.Second)
				r.Set("cache_/:time", "Wed, 16 Dec 2020 00:00:00 GMT")
				r.SetTTL("cache_/:time", 30*time.Second)
				r.Set("cache_/:type", "some-type")
				r.SetTTL("cache_/:type", 30*time.Second)
			},
			expected: &echokit.ResponseWriter{
				StatusCode: 200,
				Bytes:      []byte("\"some-response\"\n"),
				RespHeader: http.Header{
					"Cache-Control": []string{"max-age=30"},
					"Content-Type":  []string{"some-type"},
					"Expires":       []string{"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": []string{"Wed, 16 Dec 2020 00:00:00 GMT"},
				},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			testRedis, err := miniredis.Run()
			require.NoError(t, err)
			defer testRedis.Close()

			if tt.beforeFn != nil {
				tt.beforeFn(testRedis)
			}

			store := cachekit.Store{
				Client:        redis.NewClient(&redis.Options{Addr: testRedis.Addr()}),
				DefaultMaxAge: tt.defaultMaxAge,
				Prefix:        tt.prefix,
			}

			req, _ := http.NewRequest("GET", "/", nil)
			for k, v := range tt.header {
				req.Header.Add(k, v)
			}

			e := echo.New()
			w := echokit.NewResponseWriter()
			err = store.Middleware(tt.next)(e.NewContext(req, w))
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, w)
				if tt.assertFn != nil {
					tt.assertFn(t, testRedis)
				}
			}
		})
	}
}
