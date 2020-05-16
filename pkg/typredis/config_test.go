package typredis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestConfiguration(t *testing.T) {
	require.Equal(t,
		&typgo.Configuration{
			Ctor: "ctor00",
			Name: "cfgname00",
			Spec: &typredis.Config{
				Host:     "host00",
				Port:     "port00",
				Password: "pass00",
			},
		},
		typredis.Configuration(&typredis.Settings{
			Ctor:       "ctor00",
			ConfigName: "cfgname00",
			Host:       "host00",
			Port:       "port00",
			Password:   "pass00",
		}),
	)

}
