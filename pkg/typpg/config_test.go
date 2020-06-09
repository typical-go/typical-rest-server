package typpg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
)

func TestConfig(t *testing.T) {
	cfg := &typpg.Config{
		DBName:   "dbname1",
		User:     "user1",
		Password: "password1",
		Host:     "host1",
		Port:     9999,
	}

	require.Equal(t, "postgres://user1:password1@host1:9999/dbname1?sslmode=disable", typpg.Conn(cfg))
	require.Equal(t, "postgres://user1:password1@host1:9999/template1?sslmode=disable", typpg.AdminConn(cfg))
}
