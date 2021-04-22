package typdb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typdb"
)

func TestMySQL_DBTool(t *testing.T) {
	mysql := typdb.MySQL{
		Name:         "some-name",
		EnvKeys:      &typdb.EnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		DockerName:   "some-docker",
	}
	require.Equal(t, &typdb.DBTool{
		DBConn:       &typdb.MySQLConn{},
		Name:         "some-name",
		EnvKeys:      &typdb.EnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		CreateFormat: "CREATE DATABASE `%s`",
		DropFormat:   "DROP DATABASE IF EXISTS `%s`",
	}, mysql.DBTool())
}
