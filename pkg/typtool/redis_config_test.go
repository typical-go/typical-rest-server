package typtool_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

func TestRedisEnvKeys(t *testing.T) {
	os.Setenv("ABC_HOST", "some-host")
	os.Setenv("ABC_PORT", "some-port")
	os.Setenv("ABC_PASS", "some-pass")
	defer os.Clearenv()
	redisConfig := typtool.RedisEnvKeysWithPrefix("ABC")
	require.Equal(t, &typtool.RedisConfig{
		Host: "some-host",
		Port: "some-port",
		Pass: "some-pass",
	}, redisConfig.Config())
}
