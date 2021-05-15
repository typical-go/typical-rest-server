package typdb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typdb"
)

func TestPostgres_DBTool(t *testing.T) {
	pg := typdb.PostgresTool{
		Name:         "some-name",
		EnvKeys:      &typdb.EnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		DockerName:   "some-docker",
	}
	require.Equal(t, &typdb.DBTool{
		DBToolHandler: &typdb.PostgresHandler{},
		Name:          "some-name",
		EnvKeys:       &typdb.EnvKeys{},
		MigrationSrc:  "some-migr",
		SeedSrc:       "some-seed",
		CreateFormat:  "CREATE DATABASE \"%s\"",
		DropFormat:    "DROP DATABASE IF EXISTS \"%s\"",
		DockerName:    "some-docker",
	}, pg.DBTool())
}

func TestPostgres(t *testing.T) {

}
