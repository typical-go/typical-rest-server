package typdocker

import (
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// New docker module
func New() *Module {
	return &Module{}
}

// Module of docker
type Module struct{}

type docker struct {
	*typcore.Context
}

// BuildCommands is command collection to called from
func (*Module) BuildCommands(c *typcore.Context) []*cli.Command {
	d := docker{Context: c}
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				d.composeCmd(),
				d.upCmd(),
				d.downCmd(),
				d.wipeCmd(),
			},
		},
	}
}
