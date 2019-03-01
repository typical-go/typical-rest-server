package config_test

import (
	"os"
	"testing"

	"github.com/imantung/typical-go-server/app/helper/envkit"
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

func TestDetails(t *testing.T) {
	details := config.Details()
	require.Equal(t, details, []config.ConfigDetail{
		{Name: "ADDRESS", Type: "string", Required: "true"},
		{Name: "DB_NAME", Type: "string", Required: "true"},
		{Name: "DB_USER", Type: "string", Required: "true"},
		{Name: "DB_PASSWORD", Type: "string", Required: "true"},
		{Name: "DB_HOST", Type: "string", Default: "localhost"},
		{Name: "DB_PORT", Type: "int", Default: "5432"},
	})
}
