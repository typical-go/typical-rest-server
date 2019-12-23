package dbkit_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type Val struct {
	SomeField string `json:"some_field"`
}

func TestSet(t *testing.T) {
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	client := redis.NewClient(&redis.Options{Addr: testRedis.Addr()})
	err = dbkit.SetCache(client, "some-key", &Val{
		SomeField: "some-value",
	}, 3*time.Second)
	require.NoError(t, err)
	val, _ := testRedis.Get("some-key")
	require.Equal(t, "{\"some_field\":\"some-value\"}", val)
	exp := testRedis.TTL("some-key")
	require.Equal(t, 3*time.Second, exp)

}

func TestGet(t *testing.T) {
	testRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer testRedis.Close()
	client := redis.NewClient(&redis.Options{Addr: testRedis.Addr()})
	var val Val
	t.Run("WHEN no cache yet", func(t *testing.T) {
		err := dbkit.GetCache(client, "some-key", &val)
		require.EqualError(t, err, "redis: nil")
	})
	t.Run("WHEN invalid cache", func(t *testing.T) {
		testRedis.Set("some-key", `invalid-json}`)
		err := dbkit.GetCache(client, "some-key", &val)
		require.EqualError(t, err, "invalid character 'i' looking for beginning of value")
	})
	t.Run("WHEN ok", func(t *testing.T) {
		testRedis.Set("some-key", `{"some_field":"some-value"}`)
		err := dbkit.GetCache(client, "some-key", &val)
		require.NoError(t, err)
		require.Equal(t, Val{
			SomeField: "some-value",
		}, val)
	})
}
