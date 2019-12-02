package typdocker

import (
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/utility/envfile"
	"github.com/urfave/cli/v2"
)

// Module of docker
func Module() interface{} {
	return dockerModule{
		Name: "docker",
	}
}

type dockerModule struct {
	Name string
}

func (dockerModule) Commands(c *typcli.ModuleCli) []*cli.Command {
	cmd := dockerCommand{
		Context: c.Context,
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
