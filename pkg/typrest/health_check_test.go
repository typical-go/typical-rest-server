package typrest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestHealthCheck(t *testing.T) {
	testcases := []struct {
		TestName string
		typrest.HealthMap
		Expected   map[string]string
		ExpectedOk bool
	}{
		{
			HealthMap: typrest.HealthMap{
				"postgres": nil,
				"redis":    nil,
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"postgres": "OK",
				"redis":    "OK",
			},
		},
		{
			HealthMap: typrest.HealthMap{
				"postgres": errors.New("postgres-error"),
				"redis":    errors.New("redis-error"),
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"postgres": "postgres-error",
				"redis":    "redis-error",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			status, ok := tt.Status()
			require.Equal(t, tt.Expected, status)
			require.Equal(t, tt.ExpectedOk, ok)
		})
	}
}
