package typrest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestHealthCheck(t *testing.T) {
	testcases := []struct {
		TestName       string
		HealthMap      typrest.HealthMap
		Expected       bool
		ExpectedDetail map[string]string
	}{
		{
			HealthMap: typrest.HealthMap{
				"postgres": nil,
				"redis":    nil,
			},
			Expected: true,
			ExpectedDetail: map[string]string{
				"postgres": "OK",
				"redis":    "OK",
			},
		},
		{
			HealthMap: typrest.HealthMap{
				"postgres": errors.New("postgres-error"),
				"redis":    errors.New("redis-error"),
			},
			Expected: false,
			ExpectedDetail: map[string]string{
				"postgres": "postgres-error",
				"redis":    "redis-error",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			healthy, detail := typrest.HealthStatus(tt.HealthMap)
			require.Equal(t, tt.Expected, healthy)
			require.Equal(t, tt.ExpectedDetail, detail)
		})
	}
}
