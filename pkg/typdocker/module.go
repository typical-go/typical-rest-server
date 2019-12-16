package typdocker

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/utility/envfile"
	"github.com/urfave/cli/v2"
)

// Module of Docker
type Module struct{}

type docker struct {
	*typcore.Context
}

// BuildCommands is command collection to called from
func (*Module) BuildCommands(c typcore.Cli) []*cli.Command {
	d := docker{Context: c.Context()}
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Before: func(ctx *cli.Context) error {
				return envfile.Load()
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
