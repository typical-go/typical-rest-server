package typtool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

func TestMySQL_DBTool(t *testing.T) {
	mysql := typtool.MySQL{
		Name:         "some-name",
		EnvKeys:      &typtool.DBEnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		DockerName:   "some-docker",
	}
	require.Equal(t, &typtool.DBTool{
		DBConn:       &typtool.MySQLConn{},
		Name:         "some-name",
		EnvKeys:      &typtool.DBEnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		CreateFormat: "CREATE DATABASE `%s`",
		DropFormat:   "DROP DATABASE IF EXISTS `%s`",
	}, mysql.DBTool())
}
