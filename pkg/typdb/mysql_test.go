package typdb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typdb"
)

func TestMySQL_DBTool(t *testing.T) {
	envKeys := &typdb.EnvKeys{}
	mysql := typdb.MySQLTool{
		Name:         "some-name",
		EnvKeys:      envKeys,
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		DockerName:   "some-docker",
	}

	require.Equal(t, &typdb.DBTool{
		DBToolHandler: &typdb.MySQLHandler{},
		Name:          "some-name",
		EnvKeys:       envKeys,
		MigrationSrc:  "some-migr",
		SeedSrc:       "some-seed",
		CreateFormat:  "CREATE DATABASE `%s`",
		DropFormat:    "DROP DATABASE IF EXISTS `%s`",
		DockerName:    "some-docker",
	}, mysql.DBTool())
}
