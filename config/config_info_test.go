package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Yy struct {
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
}

type Xx struct {
	Host       string `envconfig:"HOST" required:"true"`
	Port       int    `envconfig:"PORT" default:"123"`
	IgnoredVar string `envconfig:"IGNORED_VAR" ignored:"true"`
	NoTagged   string

	Yy
}

func TestDetail(t *testing.T) {
	var slice []InfoDetail
	informations(&slice, &Xx{})

	require.Equal(t, slice, []InfoDetail{
		{Name: "HOST", Type: "string", Required: "true"},
		{Name: "PORT", Type: "int", Default: "123"},
		{Name: "USERNAME", Type: "string"},
		{Name: "PASSWORD", Type: "string"},
	})
}
