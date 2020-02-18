package dbkit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestCacheStore(t *testing.T) {
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()

	client := redis.NewClient(&redis.Options{Addr: testRedis.Addr()})
	store := dbkit.NewCacheStore(client)
	ctx := context.Background()

	t.Run("GIVEN no cache", func(t *testing.T) {
		t.Run("WHEN refresh failed", func(t *testing.T) {
			var b bean
			require.EqualError(t,
				store.Retrieve(ctx, "key", &b, func() (interface{}, error) {
					return nil, errors.New("some-refresh-error")
				}),
				"some-refresh-error",
			)
		})
		t.Run("WHEN marshal failed", func(t *testing.T) {
			var b bean
			require.EqualError(t,
				store.Retrieve(ctx, "key", &b, func() (interface{}, error) {
					return make(chan int), nil
				}),
				"json: unsupported type: chan int",
			)
		})
		t.Run("WHEN failed to save to redis", func(t *testing.T) {
			var b bean
			broken := dbkit.NewCacheStore(redis.NewClient(&redis.Options{Addr: "wrong-addr"}))
			require.EqualError(t,
				broken.Retrieve(ctx, "key", &b, func() (interface{}, error) {
					return &bean{Name: "new-name"}, nil
				}),
				"dial tcp: address wrong-addr: missing port in address",
			)
		})
		t.Run("", func(t *testing.T) {
			var b bean
			require.NoError(t,
				store.Retrieve(ctx, "key", &b, func() (interface{}, error) {
					return &bean{Name: "new-name"}, nil
				}),
			)
			require.Equal(t, bean{Name: "new-name"}, b)
		})
	})
	t.Run("GIVEN cache available", func(t *testing.T) {
		testRedis.Set("key", `{"name":"cached"}`)
		var b bean
		require.NoError(t,
			store.Retrieve(ctx, "key", &b, func() (interface{}, error) {
				return &bean{Name: "new-name"}, nil
			}),
		)
		require.Equal(t, bean{Name: "cached"}, b)
	})
}

type bean struct {
	Name string
}
