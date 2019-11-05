package typdocker

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicli"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
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

func (dockerModule) BuildCommand(ctx *typictx.Context) cli.Command {
	cmd := dockerCommand{
		Context: ctx,
	}
	return cli.Command{
		Name:   "docker",
		Usage:  "Docker utility",
		Before: typicli.LoadEnvFile,
		Subcommands: []cli.Command{
			{
				Name:   "compose",
				Usage:  "Generate docker-compose.yaml",
				Action: cmd.Compose,
			},
			{
				Name:  "up",
				Usage: "Create and start containers",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "no-compose",
						Usage: "Create and start containers without generate docker-compose.yaml",
					},
				},
				Action: cmd.Up,
			},
			{
				Name:   "down",
				Usage:  "Stop and remove containers, networks, images, and volumes",
				Action: cmd.Down,
			},
		},
	}
}
