package typdocker

import (
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// New docker module
func New() *Module {
	return &Module{
		Version: "3",
	}
}

// Module of docker
type Module struct {
	Version Version
}

type docker struct {
	*typcore.Context
	Version Version
}

// WithVersion to set the version
func (m *Module) WithVersion(version Version) *Module {
	m.Version = version
	return m
}

// BuildCommands is command collection to called from
func (m *Module) BuildCommands(ctx *typcore.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				m.composeCmd(ctx),
				upCmd(),
				downCmd(),
				wipeCmd(),
			},
		},
	}
}
