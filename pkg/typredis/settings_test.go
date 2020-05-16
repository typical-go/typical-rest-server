package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestInit(t *testing.T) {
	s := typredis.Init(&typredis.Settings{})
	require.Equal(t, "redis", s.Cmd)
	require.Equal(t, "REDIS", s.ConfigName)
	require.Equal(t, "localhost", s.Host)
	require.Equal(t, "6379", s.Port)
	require.Equal(t, "redis", s.DockerName)
	require.Equal(t, "redis:4.0.5-alpine", s.DockerImage)
}
