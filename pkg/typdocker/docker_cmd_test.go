package typdocker_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli/v2"
)

func TestCmdUp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{
			{CommandLine: []string{"docker-compose", "up", "--remove-orphans", "-d"}},
		})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdUp(&typgo.BuildSys{})
		cliCtx := cli.NewContext(nil, &flag.FlagSet{}, nil)
		require.NoError(t, command.Action(cliCtx))
	})
	t.Run("with wipe", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{
			{CommandLine: []string{"docker", "ps", "-q"}},
			{CommandLine: []string{"docker-compose", "up", "--remove-orphans", "-d"}},
		})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdUp(&typgo.BuildSys{})
		flagSet := &flag.FlagSet{}
		flagSet.Bool("wipe", true, "")

		require.NoError(t, command.Action(cli.NewContext(nil, flagSet, nil)))
	})

	t.Run("with wipe error", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdUp(&typgo.BuildSys{})
		flagSet := &flag.FlagSet{}
		flagSet.Bool("wipe", true, "")

		err := command.Action(cli.NewContext(nil, flagSet, nil))
		require.EqualError(t, err, "Docker-ID: execkit-mock: no run expectation for [docker ps -q]")
	})
}

func TestCmdWipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{
			{CommandLine: []string{"docker", "ps", "-q"}, OutputBytes: []byte("pid-1\npid-2")},
			{CommandLine: []string{"docker", "kill", "pid-1"}},
			{CommandLine: []string{"docker", "kill", "pid-2"}},
		})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdWipe(&typgo.BuildSys{})

		require.NoError(t, command.Action(cli.NewContext(nil, &flag.FlagSet{}, nil)))
	})

	t.Run("when ps error", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdWipe(&typgo.BuildSys{})

		err := command.Action(cli.NewContext(nil, &flag.FlagSet{}, nil))
		require.EqualError(t, err, "Docker-ID: execkit-mock: no run expectation for [docker ps -q]")
	})

	t.Run("when kill error", func(t *testing.T) {
		unpatch := execkit.Patch([]*execkit.RunExpectation{
			{CommandLine: []string{"docker", "ps", "-q"}, OutputBytes: []byte("pid-1\npid-2")},
		})
		defer unpatch(t)
		cmd := &typdocker.DockerCmd{}
		command := cmd.CmdWipe(&typgo.BuildSys{})

		err := command.Action(cli.NewContext(nil, &flag.FlagSet{}, nil))
		require.EqualError(t, err, "Fail to kill #pid-1: execkit-mock: no run expectation for [docker kill pid-1]")
	})

}

func TestCmdDown(t *testing.T) {
	unpatch := execkit.Patch([]*execkit.RunExpectation{
		{CommandLine: []string{"docker-compose", "down", "-v"}},
	})
	defer unpatch(t)
	cmd := &typdocker.DockerCmd{}
	command := cmd.CmdDown(&typgo.BuildSys{})

	require.NoError(t, command.Action(cli.NewContext(nil, &flag.FlagSet{}, nil)))
}
