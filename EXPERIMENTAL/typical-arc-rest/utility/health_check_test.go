package utility_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typical-arc-rest/utility"
)

func TestHealthCheck_SetStatusOK(t *testing.T) {
	healthcheck := utility.NewHealthCheck().
		Add("component_1", nil).
		Add("component_2", fmt.Errorf("some error"))

	require.Equal(t, healthcheck["component_1"], "OK")
	require.Equal(t, healthcheck["component_2"], "some error")
}

func TestHealthCheck_NotOK(t *testing.T) {
	testcases := []struct {
		HealthCheck utility.HealthCheck
		NotOK       bool
	}{
		{
			utility.NewHealthCheck().
				Add("component_1", nil).
				Add("component_2", nil),
			false,
		},
		{
			utility.NewHealthCheck().
				Add("component_1", nil).
				Add("component_2", fmt.Errorf("some error")),
			true,
		},
		{
			utility.NewHealthCheck().
				Add("component_1", fmt.Errorf("some error")).
				Add("component_2", fmt.Errorf("some error")),
			true,
		},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.HealthCheck.NotOK(), tt.NotOK)
	}
}
