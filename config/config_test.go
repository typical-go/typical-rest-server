package config_test

import (
	"os"
	"testing"

	"github.com/typical-go/typical-rest-server/app/helper/envkit"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/stretchr/testify/require"
)

var (
	all_good_env = map[string]string{
		"TEST_ADDRESS":     ":99999",
		"TEST_PG_DBNAME":   "some-dbname",
		"TEST_PG_USER":     "some-dbuser",
		"TEST_PG_PASSWORD": "some-dbpassword",
		"TEST_PG_HOST":     "some-dbhost",
		"TEST_PG_PORT":     "88888",
	}
)

func init() {
	config.Prefix = "TEST"
}

func TestLoadConfig(t *testing.T) {
	envkit.Set(all_good_env)
	defer os.Clearenv()
	conf, err := config.LoadConfig()
	require.NoError(t, err)
	require.Equal(t, conf.Address, ":99999")
	require.Equal(t, conf.Postgres.DbName, "some-dbname")
	require.Equal(t, conf.Postgres.Password, "some-dbpassword")
	require.Equal(t, conf.Postgres.Host, "some-dbhost")
	require.Equal(t, conf.Postgres.Port, 88888)

}

func TestDetails(t *testing.T) {
	details := config.Informations()
	require.Equal(t, details, []config.InfoDetail{
		{Name: "ADDRESS", Type: "string", Required: "true"},
		{Name: "PG_DBNAME", Type: "string", Required: "true"},
		{Name: "PG_USER", Type: "string", Required: "true"},
		{Name: "PG_PASSWORD", Type: "string", Required: "true"},
		{Name: "PG_HOST", Type: "string", Default: "localhost"},
		{Name: "PG_PORT", Type: "int", Default: "5432"},
	})
}
