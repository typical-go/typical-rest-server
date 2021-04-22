package typdocker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

func TestCmdWipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
			{CommandLine: "docker kill pid-1"},
			{CommandLine: "docker kill pid-2"},
		})(t)

		tool := &typdocker.DockerTool{}
		require.NoError(t, tool.DockerWipe(c))
	})

	t.Run("when ps error", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{})(t)

		tool := &typdocker.DockerTool{}
		err := tool.DockerWipe(c)
		require.EqualError(t, err, "DockerTool-ID: typgo-mock: no run expectation for \"docker ps -q\"")
	})

	t.Run("when kill error", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
		})(t)

		tool := &typdocker.DockerTool{}
		err := tool.DockerWipe(c)
		require.EqualError(t, err, "fail to kill #pid-1: typgo-mock: no run expectation for \"docker kill pid-1\"")
	})

}
