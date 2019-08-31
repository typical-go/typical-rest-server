package typredis_test

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
	"github.com/tiket/TIX-COMMON-GO/tixredis"
)

func TestConnect(t *testing.T) {
	dummyServer, err := miniredis.Run()
	require.NoError(t, err)
	defer dummyServer.Close()
	t.Run("GIVEN bad config", func(t *testing.T) {
		_, err := tixredis.Connect(&tixredis.Config{
			Host: "badserver",
			Port: "1",
		})
		require.EqualError(t, err, "dial tcp: lookup badserver: no such host")
	})
	t.Run("GIVEN good config", func(t *testing.T) {
		dummyServer.Set("hello", "world")
		client, err := tixredis.Connect(&tixredis.Config{
			Host: dummyServer.Host(),
			Port: dummyServer.Port(),
		})
		require.NoError(t, err)
		require.Equal(t, "world", client.Get("hello").Val())
	})
}
