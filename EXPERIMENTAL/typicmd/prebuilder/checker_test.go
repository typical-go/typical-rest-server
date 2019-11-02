package prebuilder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChecker_CheckBuildTool(t *testing.T) {
	testcases := []struct {
		checker        checker
		checkBuildTool bool
	}{
		{checker{mockTarget: true}, true},
		{checker{constructor: true}, true},
		{checker{configuration: true}, true},
		{checker{testTarget: true}, true},
		{checker{buildToolBinary: true}, true},
		{checker{contextChecksum: true}, true},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.checkBuildTool, tt.checker.checkBuildTool())
	}
}
