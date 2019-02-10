package config_test

import (
	"os"
	"testing"

	"github.com/imantung/go-helper/envkit"
	"github.com/imantung/typical-go-server/config"
	"github.com/stretchr/testify/require"
)

var (
	all_good_env = map[string]string{
		"TEST_ADDRESS":     ":99999",
		"TEST_DB_NAME":     "some-dbname",
		"TEST_DB_USER":     "some-dbuser",
		"TEST_DB_PASSWORD": "some-dbpassword",
		"TEST_DB_HOST":     "some-dbhost",
		"TEST_DB_PORT":     "88888",
	}
)

func init() {
	config.Prefix = "TEST"
}

func TestConfig(t *testing.T) {
	envkit.Set(all_good_env)
	defer os.Clearenv()
	conf, err := config.NewConfig()
	require.NoError(t, err)
	require.Equal(t, conf.Address, ":99999")
}

func TestConfig_WhenNoConfig(t *testing.T) {

	t.Run("ReturnDefaultValue", func(t *testing.T) {
		envkit.Set(all_good_env)
		os.Unsetenv("TEST_DB_HOST")
		os.Unsetenv("TEST_DB_PORT")
		defer os.Clearenv()
		conf, err := config.NewConfig()
		require.NoError(t, err)
		require.Equal(t, conf.DbHost, "localhost")
		require.Equal(t, conf.DbPort, 5432)

	})

	t.Run("ReturnError", func(t *testing.T) {
		defer os.Clearenv()
		_, err := config.NewConfig()
		require.EqualError(t, err, "required key TEST_ADDRESS missing value")

		os.Setenv("TEST_ADDRESS", "some-address")
		_, err = config.NewConfig()
		require.EqualError(t, err, "required key TEST_DB_NAME missing value")

		os.Setenv("TEST_DB_NAME", "some-dbname")
		_, err = config.NewConfig()
		require.EqualError(t, err, "required key TEST_DB_USER missing value")

		os.Setenv("TEST_DB_USER", "some-dbuser")
		_, err = config.NewConfig()
		require.EqualError(t, err, "required key TEST_DB_PASSWORD missing value")

		os.Setenv("TEST_DB_PASSWORD", "some-dbpassword")
		_, err = config.NewConfig()
		require.NoError(t, err)
	})
}
