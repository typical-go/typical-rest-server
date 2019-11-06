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
		{checker{}, false},
		{checker{mockTarget: true}, true},
		{checker{configuration: true}, true},
		{checker{testTarget: true}, true},
		{checker{buildToolBinary: true}, true},
		{checker{contextChecksum: true}, true},
		{checker{buildCommands: true}, true},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.checkBuildTool, tt.checker.checkBuildTool())
	}
}

func TestChecker_CheckReadme(t *testing.T) {
	testcases := []struct {
		checker     checker
		checkReadme bool
	}{
		{checker{}, false},
		{checker{configuration: true}, true},
		{checker{buildCommands: true}, true},
		{checker{readmeFile: true}, true},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.checkReadme, tt.checker.checkReadme())
	}
}
