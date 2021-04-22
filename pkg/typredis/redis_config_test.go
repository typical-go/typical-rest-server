package typredis_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestEnvKeys(t *testing.T) {
	os.Setenv("ABC_HOST", "some-host")
	os.Setenv("ABC_PORT", "some-port")
	os.Setenv("ABC_PASS", "some-pass")
	defer os.Clearenv()
	Config := typredis.EnvKeysWithPrefix("ABC")
	require.Equal(t, &typredis.Config{
		Host: "some-host",
		Port: "some-port",
		Pass: "some-pass",
	}, Config.Config())
}
