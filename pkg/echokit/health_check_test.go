package echokit_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

type (
	healtCheck struct {
		testName        string
		healthCheck     echokit.HealthCheck
		expectedStatus  int
		expectedMessage map[string]string
	}
)

func TestHealthCheck(t *testing.T) {
	testcases := []healtCheck{
		{
			healthCheck: echokit.HealthCheck{
				"postgres": func() error { return nil },
				"redis":    func() error { return nil },
			},
			expectedStatus: 200,
			expectedMessage: map[string]string{
				"postgres": "OK",
				"redis":    "OK",
			},
		},
		{
			healthCheck: echokit.HealthCheck{
				"postgres": func() error { return errors.New("postgres-error") },
				"redis":    func() error { return errors.New("redis-error") },
			},
			expectedStatus: 503,
			expectedMessage: map[string]string{
				"postgres": "postgres-error",
				"redis":    "redis-error",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			status, message := tt.healthCheck.Result()
			require.Equal(t, tt.expectedStatus, status)
			require.Equal(t, tt.expectedMessage, message)
		})
	}
}
