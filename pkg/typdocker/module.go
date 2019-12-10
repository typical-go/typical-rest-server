package typdocker

import (
	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-go/pkg/utility/envfile"
	"github.com/urfave/cli/v2"
)

// Module of Docker
type Module struct{}

// BuildCommands is command collection to called from
func (*Module) BuildCommands(c typobj.Cli) []*cli.Command {
	cmd := dockerCommand{
		context: c.(context),
	}
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Before: func(ctx *cli.Context) error {
				return envfile.Load()
			},
			Subcommands: []*cli.Command{
				{
					Name:   "compose",
					Usage:  "Generate docker-compose.yaml",
					Action: cmd.Compose,
				},
				{
					Name:   "up",
					Usage:  "Spin up docker containers",
					Action: cmd.Up,
				},
				{
					Name:   "down",
					Usage:  "Take down all docker containers",
					Action: cmd.Down,
				},
			},
		},
	}
}
