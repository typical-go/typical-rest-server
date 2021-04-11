package typtool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

func TestPostgres_DBTool(t *testing.T) {
	mysql := typtool.Postgres{
		Name:         "some-name",
		EnvKeys:      &typtool.DBEnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		DockerName:   "some-docker",
	}
	require.Equal(t, &typtool.DBTool{
		DBConn:       &typtool.PGConn{},
		Name:         "some-name",
		EnvKeys:      &typtool.DBEnvKeys{},
		MigrationSrc: "some-migr",
		SeedSrc:      "some-seed",
		CreateFormat: "CREATE DATABASE \"%s\"",
		DropFormat:   "DROP DATABASE IF EXISTS \"%s\"",
	}, mysql.DBTool())
}

func TestPostgres(t *testing.T) {

}
