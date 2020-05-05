package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestSettingFunctions(t *testing.T) {
	require.Equal(t, "postgres", typpostgres.GetDockerName(&typpostgres.Setting{}))
	require.Equal(t, "postgres", typpostgres.GetDockerImage(&typpostgres.Setting{}))
	require.Equal(t, "PG", typpostgres.GetConfigName(&typpostgres.Setting{}))
	require.Equal(t, "sample", typpostgres.GetDBName(&typpostgres.Setting{}))
	require.Equal(t, "postgres", typpostgres.GetUser(&typpostgres.Setting{}))
	require.Equal(t, "pgpass", typpostgres.GetPassword(&typpostgres.Setting{}))
	require.Equal(t, "localhost", typpostgres.GetHost(&typpostgres.Setting{}))
	require.Equal(t, 5432, typpostgres.GetPort(&typpostgres.Setting{}))
	require.Equal(t, "scripts/db/migration", typpostgres.GetMigrationSrc(&typpostgres.Setting{}))
	require.Equal(t, "scripts/db/seed", typpostgres.GetSeedSrc(&typpostgres.Setting{}))

	require.Equal(t,
		"docker-name-00",
		typpostgres.GetDockerName(&typpostgres.Setting{DockerName: "docker-name-00"}),
	)
	require.Equal(t,
		"docker-image-00",
		typpostgres.GetDockerImage(&typpostgres.Setting{DockerImage: "docker-image-00"}),
	)

	require.Equal(t,
		"config-name-00",
		typpostgres.GetConfigName(&typpostgres.Setting{ConfigName: "config-name-00"}),
	)

	require.Equal(t,
		"db-name-00",
		typpostgres.GetDBName(&typpostgres.Setting{DBName: "db-name-00"}),
	)

	require.Equal(t,
		"user-00",
		typpostgres.GetUser(&typpostgres.Setting{User: "user-00"}),
	)

	require.Equal(t,
		"password-00",
		typpostgres.GetPassword(&typpostgres.Setting{Password: "password-00"}),
	)

	require.Equal(t,
		"host-00",
		typpostgres.GetHost(&typpostgres.Setting{Host: "host-00"}),
	)

	require.Equal(t,
		9999,
		typpostgres.GetPort(&typpostgres.Setting{Port: 9999}),
	)
}
