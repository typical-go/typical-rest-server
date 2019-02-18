package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetail(t *testing.T) {
	var conf struct {
		Host       string `envconfig:"HOST" required:"true"`
		Port       int    `envconfig:"PORT" default:"123"`
		IgnoredVar string `envconfig:"IGNORED_VAR" ignored:"true"`
		NoTagged   string
	}

	details := details(&conf)

	require.Equal(t, details, []ConfigDetail{
		{Name: "HOST", Type: "string", Required: "true"},
		{Name: "PORT", Type: "int", Default: "123"},
	})
}
