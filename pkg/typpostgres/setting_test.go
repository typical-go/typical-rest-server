package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestSettingFunctions(t *testing.T) {
	require.Equal(t, "postgres", typpostgres.DockerName(&typpostgres.Setting{}))
	require.Equal(t, "postgres", typpostgres.DockerImage(&typpostgres.Setting{}))
	require.Equal(t, "PG", typpostgres.ConfigName(&typpostgres.Setting{}))
	require.Equal(t, "sample", typpostgres.DBName(&typpostgres.Setting{}))
	require.Equal(t, "postgres", typpostgres.User(&typpostgres.Setting{}))
	require.Equal(t, "pgpass", typpostgres.Password(&typpostgres.Setting{}))
	require.Equal(t, "localhost", typpostgres.Host(&typpostgres.Setting{}))
	require.Equal(t, 5432, typpostgres.Port(&typpostgres.Setting{}))
	require.Equal(t, "scripts/db/migration", typpostgres.MigrationSrc(&typpostgres.Setting{}))
	require.Equal(t, "scripts/db/seed", typpostgres.SeedSrc(&typpostgres.Setting{}))

	require.Equal(t,
		"docker-name-00",
		typpostgres.DockerName(&typpostgres.Setting{DockerName: "docker-name-00"}),
	)
	require.Equal(t,
		"docker-image-00",
		typpostgres.DockerImage(&typpostgres.Setting{DockerImage: "docker-image-00"}),
	)

	require.Equal(t,
		"config-name-00",
		typpostgres.ConfigName(&typpostgres.Setting{ConfigName: "config-name-00"}),
	)

	require.Equal(t,
		"db-name-00",
		typpostgres.DBName(&typpostgres.Setting{DBName: "db-name-00"}),
	)

	require.Equal(t,
		"user-00",
		typpostgres.User(&typpostgres.Setting{User: "user-00"}),
	)

	require.Equal(t,
		"password-00",
		typpostgres.Password(&typpostgres.Setting{Password: "password-00"}),
	)

	require.Equal(t,
		"host-00",
		typpostgres.Host(&typpostgres.Setting{Host: "host-00"}),
	)

	require.Equal(t,
		9999,
		typpostgres.Port(&typpostgres.Setting{Port: 9999}),
	)

}
