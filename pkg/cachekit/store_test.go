package cachekit_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/typical-go/typical-rest-server/pkg/echotest"
)

func TestStore_Middleware(t *testing.T) {
	defer monkey.Patch(time.Now, func() time.Time {
		return time.Date(2020, time.December, 16, 0, 0, 0, 0, time.UTC)
	}).Unpatch()

	testcases := []struct {
		testName      string
		next          echo.HandlerFunc
		defaultMaxAge time.Duration
		prefixKey     string
		header        map[string]string
		beforeFn      func(*miniredis.Miniredis)
		assertFn      func(*testing.T, *miniredis.Miniredis)
		expected      echotest.Response
		expectedErr   string
	}{
		{
			testName: "no cache available",
			next: func(ec echo.Context) error {
				return ec.JSON(200, "some-response")
			},
			defaultMaxAge: 30 * time.Second,
			prefixKey:     "cache_",
			expected: echotest.Response{
				Code: 200,
				Body: "\"some-response\"\n",
				Header: http.Header{
					"Cache-Control": {"max-age=30"},
					"Content-Type":  {"application/json; charset=UTF-8"},
					"Expires":       {"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": {"Wed, 16 Dec 2020 00:00:00 GMT"},
				},
			},

			assertFn: func(t *testing.T, r *miniredis.Miniredis) {
				data, _ := r.Get("cache_/")
				require.Equal(t, "{\"LastModified\":\"Wed, 16 Dec 2020 00:00:00 GMT\",\"Body\":\"InNvbWUtcmVzcG9uc2UiCg==\",\"Head\":{\"StatusCode\":200,\"Header\":{\"Content-Type\":[\"application/json; charset=UTF-8\"]}}}", data)
				require.Equal(t, 30*time.Second, r.TTL("cache_/"))
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
			prefixKey:     "cache_",
			expected: echotest.Response{
				Code: 200,
				Body: "\"some-response\"\n",
				Header: http.Header{
					"Cache-Control": {"max-age=30"},
					"Content-Type":  {"application/json; charset=UTF-8"},
					"X-Total-Count": {"6"},
					"Expires":       {"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": {"Wed, 16 Dec 2020 00:00:00 GMT"},
				},
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "{\"LastModified\":\"Wed, 16 Dec 2020 00:00:00 GMT\",\"Body\":\"InNvbWUtcmVzcG9uc2UiCg==\",\"Head\":{\"StatusCode\":200,\"Header\":{\"Content-Type\":[\"application/json; charset=UTF-8\"],\"X-Total-Count\":[\"6\"]}}}")
				r.SetTTL("cache_/", 30*time.Second)
			},
		},
		{
			testName:      "if cache not modified since",
			defaultMaxAge: 30 * time.Second,
			prefixKey:     "cache_",
			header: map[string]string{
				"If-Modified-Since": "Wed, 16 Dec 2020 00:00:05 GMT",
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "{\"LastModified\":\"Wed, 16 Dec 2020 00:00:00 GMT\",\"Body\":\"InNvbWUtcmVzcG9uc2UiCg==\",\"Head\":{\"StatusCode\":200,\"Header\":{}}}")
				r.SetTTL("cache_/", 30*time.Second)
			},
			expectedErr: "code=304, message=Not Modified",
		},
		{
			testName:      "if cache modified since",
			defaultMaxAge: 30 * time.Second,
			prefixKey:     "cache_",
			next: func(ec echo.Context) error {
				return ec.JSON(200, "some-response")
			},
			header: map[string]string{
				"If-Modified-Since": "Wed, 15 Dec 2020 23:59:00 GMT",
			},
			beforeFn: func(r *miniredis.Miniredis) {
				r.Set("cache_/", "{\"LastModified\":\"Wed, 16 Dec 2020 00:00:00 GMT\",\"Body\":\"InNvbWUtcmVzcG9uc2UiCg==\",\"Head\":{\"StatusCode\":200,\"Header\":{\"Content-Type\":[\"application/json; charset=UTF-8\"]}}}")
				r.SetTTL("cache_/", 30*time.Second)
			},
			expected: echotest.Response{
				Code: 200,
				Body: "\"some-response\"\n",
				Header: http.Header{
					"Cache-Control": {"max-age=30"},
					"Content-Type":  {"application/json; charset=UTF-8"},
					"Expires":       {"Wed, 16 Dec 2020 00:00:30 GMT"},
					"Last-Modified": {"Wed, 16 Dec 2020 00:00:00 GMT"},
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
				PrefixKey:     tt.prefixKey,
			}

			req, _ := http.NewRequest("GET", "/", nil)
			for k, v := range tt.header {
				req.Header.Add(k, v)
			}

			e := echo.New()
			rec := httptest.NewRecorder()
			err = store.Middleware(tt.next)(e.NewContext(req, rec))
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				echotest.EqualResp(t, tt.expected, rec)
				if tt.assertFn != nil {
					tt.assertFn(t, testRedis)
				}
			}
		})
	}
}
