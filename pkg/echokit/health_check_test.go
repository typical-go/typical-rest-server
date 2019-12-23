package echokit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestHealthCheck_SetStatusOK(t *testing.T) {
	healthcheck := echokit.NewHealthCheck().
		Add("component_1", nil).
		Add("component_2", fmt.Errorf("some error"))

	require.Equal(t, healthcheck["component_1"], "OK")
	require.Equal(t, healthcheck["component_2"], "some error")
}

func TestHealthCheck_NotOK(t *testing.T) {
	testcases := []struct {
		HealthCheck echokit.HealthCheck
		NotOK       bool
	}{
		{
			echokit.NewHealthCheck().
				Add("component_1", nil).
				Add("component_2", nil),
			false,
		},
		{
			echokit.NewHealthCheck().
				Add("component_1", nil).
				Add("component_2", fmt.Errorf("some error")),
			true,
		},
		{
			echokit.NewHealthCheck().
				Add("component_1", fmt.Errorf("some error")).
				Add("component_2", fmt.Errorf("some error")),
			true,
		},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.HealthCheck.NotOK(), tt.NotOK)
	}
}
