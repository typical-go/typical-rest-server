package typdocker_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli/v2"
)

func TestCmdUp(t *testing.T) {
	c := &typgo.Context{}
	defer c.PatchBash([]*typgo.MockBash{
		{CommandLine: "docker-compose up --remove-orphans -d"},
	})(t)
	err := typdocker.DockerUp(c)
	require.NoError(t, err)
}

func TestCmdUp_WithPipe(t *testing.T) {
	flagSet := &flag.FlagSet{}
	flagSet.Bool("wipe", true, "")
	cc := cli.NewContext(nil, flagSet, nil)
	cc.Command = &cli.Command{}

	c := &typgo.Context{
		Context:    cc,
		Descriptor: &typgo.Descriptor{},
	}
	defer c.PatchBash([]*typgo.MockBash{
		{CommandLine: "docker ps -q"},
		{CommandLine: "docker-compose up --remove-orphans -d"},
	})(t)

	err := typdocker.DockerUp(c)

	require.NoError(t, err)
}

func TestCmdUp_WithPipeError(t *testing.T) {
	flagSet := &flag.FlagSet{}
	flagSet.Bool("wipe", true, "")

	cc := cli.NewContext(nil, flagSet, nil)
	cc.Command = &cli.Command{Name: "dummy"}

	c := &typgo.Context{
		Context:    cc,
		Descriptor: &typgo.Descriptor{},
	}
	defer c.PatchBash([]*typgo.MockBash{})(t)

	err := typdocker.DockerUp(c)

	require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
}

func TestCmdWipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
			{CommandLine: "docker kill pid-1"},
			{CommandLine: "docker kill pid-2"},
		})(t)
		require.NoError(t, typdocker.DockerWipe(c))
	})

	t.Run("when ps error", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{})(t)
		err := typdocker.DockerWipe(c)
		require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
	})

	t.Run("when kill error", func(t *testing.T) {
		c := &typgo.Context{}
		defer c.PatchBash([]*typgo.MockBash{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
		})(t)
		err := typdocker.DockerWipe(c)
		require.EqualError(t, err, "Fail to kill #pid-1: typgo-mock: no run expectation for \"docker kill pid-1\"")
	})

}

func TestCmdDown(t *testing.T) {
	c := &typgo.Context{}
	defer c.PatchBash([]*typgo.MockBash{
		{CommandLine: "docker-compose down -v"},
	})(t)
	require.NoError(t, typdocker.DockerDown(c))
}
