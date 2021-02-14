package typdocker_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli/v2"
)

func TestCmdUp(t *testing.T) {
	defer typgo.PatchBash([]*typgo.RunExpectation{
		{CommandLine: "docker-compose up --remove-orphans -d"},
	})(t)

	c, _ := typgo.DummyContext()
	err := typdocker.DockerUp(c)
	require.NoError(t, err)
}

func TestCmdUp_WithPipe(t *testing.T) {
	defer typgo.PatchBash([]*typgo.RunExpectation{
		{CommandLine: "docker ps -q"},
		{CommandLine: "docker-compose up --remove-orphans -d"},
	})(t)

	flagSet := &flag.FlagSet{}
	flagSet.Bool("wipe", true, "")
	c := cli.NewContext(nil, flagSet, nil)
	c.Command = &cli.Command{}

	err := typdocker.DockerUp(&typgo.Context{
		Context:    c,
		Descriptor: &typgo.Descriptor{},
		Stdout:     &strings.Builder{},
	})

	require.NoError(t, err)
}

func TestCmdUp_WithPipeError(t *testing.T) {
	defer typgo.PatchBash([]*typgo.RunExpectation{})(t)

	flagSet := &flag.FlagSet{}
	flagSet.Bool("wipe", true, "")

	c := cli.NewContext(nil, flagSet, nil)
	c.Command = &cli.Command{Name: "dummy"}

	err := typdocker.DockerUp(&typgo.Context{
		Context:    c,
		Descriptor: &typgo.Descriptor{},
		Stdout:     &strings.Builder{},
	})

	require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
}

func TestCmdWipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
			{CommandLine: "docker kill pid-1"},
			{CommandLine: "docker kill pid-2"},
		})(t)

		c, _ := typgo.DummyContext()
		require.NoError(t, typdocker.DockerWipe(c))
	})

	t.Run("when ps error", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{})(t)

		c, _ := typgo.DummyContext()
		err := typdocker.DockerWipe(c)
		require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
	})

	t.Run("when kill error", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
		})(t)

		c, _ := typgo.DummyContext()
		err := typdocker.DockerWipe(c)
		require.EqualError(t, err, "Fail to kill #pid-1: typgo-mock: no run expectation for \"docker kill pid-1\"")
	})

}

func TestCmdDown(t *testing.T) {
	defer typgo.PatchBash([]*typgo.RunExpectation{
		{CommandLine: "docker-compose down -v"},
	})(t)

	c, _ := typgo.DummyContext()
	require.NoError(t, typdocker.DockerDown(c))
}
