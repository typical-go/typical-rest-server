package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestSettingFunctions(t *testing.T) {
	require.Equal(t, "postgres", typpostgres.GetDockerName(&typpostgres.Settings{}))
	require.Equal(t, "postgres", typpostgres.GetDockerImage(&typpostgres.Settings{}))
	require.Equal(t, "PG", typpostgres.GetConfigName(&typpostgres.Settings{}))
	require.Equal(t, "sample", typpostgres.GetDBName(&typpostgres.Settings{}))
	require.Equal(t, "postgres", typpostgres.GetUser(&typpostgres.Settings{}))
	require.Equal(t, "pgpass", typpostgres.GetPassword(&typpostgres.Settings{}))
	require.Equal(t, "localhost", typpostgres.GetHost(&typpostgres.Settings{}))
	require.Equal(t, 5432, typpostgres.GetPort(&typpostgres.Settings{}))
	require.Equal(t, "scripts/db/migration", typpostgres.GetMigrationSrc(&typpostgres.Settings{}))
	require.Equal(t, "scripts/db/seed", typpostgres.GetSeedSrc(&typpostgres.Settings{}))

	require.Equal(t,
		"docker-name-00",
		typpostgres.GetDockerName(&typpostgres.Settings{DockerName: "docker-name-00"}),
	)
	require.Equal(t,
		"docker-image-00",
		typpostgres.GetDockerImage(&typpostgres.Settings{DockerImage: "docker-image-00"}),
	)

	require.Equal(t,
		"config-name-00",
		typpostgres.GetConfigName(&typpostgres.Settings{ConfigName: "config-name-00"}),
	)

	require.Equal(t,
		"db-name-00",
		typpostgres.GetDBName(&typpostgres.Settings{DBName: "db-name-00"}),
	)

	require.Equal(t,
		"user-00",
		typpostgres.GetUser(&typpostgres.Settings{User: "user-00"}),
	)

	require.Equal(t,
		"password-00",
		typpostgres.GetPassword(&typpostgres.Settings{Password: "password-00"}),
	)

	require.Equal(t,
		"host-00",
		typpostgres.GetHost(&typpostgres.Settings{Host: "host-00"}),
	)

	require.Equal(t,
		9999,
		typpostgres.GetPort(&typpostgres.Settings{Port: 9999}),
	)
}
