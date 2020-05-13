package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestConfig(t *testing.T) {
	cfg := &typpostgres.Config{
		DBName:   "dbname1",
		User:     "user1",
		Password: "password1",
		Host:     "host1",
		Port:     9999,
	}

	require.Equal(t,
		"postgres://user1:password1@host1:9999/dbname1?sslmode=disable",
		cfg.ConnStr(),
	)
	require.Equal(t,
		"postgres://user1:password1@host1:9999/template1?sslmode=disable",
		cfg.Admin().ConnStr(),
	)
}
