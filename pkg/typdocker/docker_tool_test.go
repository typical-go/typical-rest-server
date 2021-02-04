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
	t.Run("success", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker-compose up --remove-orphans -d"},
		})(t)

		err := typdocker.DockerUp(&typgo.Context{
			Context: cli.NewContext(nil, &flag.FlagSet{}, nil),
		})
		require.NoError(t, err)
	})
	t.Run("with wipe", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker ps -q"},
			{CommandLine: "docker-compose up --remove-orphans -d"},
		})(t)

		flagSet := &flag.FlagSet{}
		flagSet.Bool("wipe", true, "")

		err := typdocker.DockerUp(&typgo.Context{
			Context: cli.NewContext(nil, flagSet, nil),
		})

		require.NoError(t, err)
	})

	t.Run("with wipe error", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{})(t)

		flagSet := &flag.FlagSet{}
		flagSet.Bool("wipe", true, "")

		err := typdocker.DockerUp(&typgo.Context{
			Context: cli.NewContext(nil, flagSet, nil),
		})

		require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
	})
}

func TestCmdWipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
			{CommandLine: "docker kill pid-1"},
			{CommandLine: "docker kill pid-2"},
		})(t)

		require.NoError(t, typdocker.DockerWipe(&typgo.Context{
			Context: cli.NewContext(nil, &flag.FlagSet{}, nil),
		}))
	})

	t.Run("when ps error", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{})(t)

		err := typdocker.DockerWipe(&typgo.Context{
			Context: cli.NewContext(nil, &flag.FlagSet{}, nil),
		})
		require.EqualError(t, err, "Docker-ID: typgo-mock: no run expectation for \"docker ps -q\"")
	})

	t.Run("when kill error", func(t *testing.T) {
		defer typgo.PatchBash([]*typgo.RunExpectation{
			{CommandLine: "docker ps -q", OutputBytes: []byte("pid-1\npid-2")},
		})(t)

		err := typdocker.DockerWipe(&typgo.Context{
			Context: cli.NewContext(nil, &flag.FlagSet{}, nil),
		})
		require.EqualError(t, err, "Fail to kill #pid-1: typgo-mock: no run expectation for \"docker kill pid-1\"")
	})

}

func TestCmdDown(t *testing.T) {
	defer typgo.PatchBash([]*typgo.RunExpectation{
		{CommandLine: "docker-compose down -v"},
	})(t)

	require.NoError(t, typdocker.DockerDown(&typgo.Context{
		Context: cli.NewContext(nil, &flag.FlagSet{}, nil),
	}))
}
