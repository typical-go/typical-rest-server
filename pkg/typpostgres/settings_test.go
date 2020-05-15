package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestSettingFunctions(t *testing.T) {
	s := typpostgres.Init(&typpostgres.Settings{})

	require.Equal(t, "pg", s.UtilityCmd)
	require.Equal(t, "postgres", s.DockerName)
	require.Equal(t, "postgres", s.DockerImage)
	require.Equal(t, "PG", s.ConfigName)
	require.Equal(t, "sample", s.DBName)
	require.Equal(t, "postgres", s.User)
	require.Equal(t, "pgpass", s.Password)
	require.Equal(t, "localhost", s.Host)
	require.Equal(t, 5432, s.Port)

}
