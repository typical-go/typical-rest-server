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
	}
}
