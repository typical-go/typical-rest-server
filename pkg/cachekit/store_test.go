package cachekit_test

import (
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/typical-go/typical-rest-server/pkg/echokit"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

func TestStore_Middleware(t *testing.T) {
	// monkey patch time.Now
	defer monkey.Patch(time.Now, func() time.Time {
		return time.Date(2020, time.December, 16, 0, 0, 0, 0, time.UTC)
	}).Unpatch()

	testcases := []struct {
		testName      string
		next          echo.HandlerFunc
		defaultMaxAge time.Duration
		header        map[string]string
		beforeFn      func(*miniredis.Miniredis)
		assertFn      func(*testing.T, *miniredis.Miniredis)
		expected      *echokit.ResponseWriter
		expectedErr   string
	}{
		{
			next: func(ec echo.Context) error {
				return ec.JSON(200, "some-response")
			},
			defaultMaxAge: 30 * time.Second,
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
				data, _ := r.Get("/")
				lastModified, _ := r.Get("/:time")

				require.Equal(t, "\"some-response\"\n", data)
				require.Equal(t, "Wed, 16 Dec 2020 00:00:00 GMT", lastModified)

				require.Equal(t, 30*time.Second, r.TTL("/"))
				require.Equal(t, 30*time.Second, r.TTL("/:time"))
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
