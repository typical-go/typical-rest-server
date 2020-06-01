package typical

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
	"github.com/urfave/cli/v2"
)

type redisDocker struct{}

//
// util
//

func redisUtil(c *typgo.BuildCli) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis utility",
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Action:  c.ActionFn("REDIS", redisConsole),
				},
			},
		},
	}
}

func redisConsole(c *typgo.Context) (err error) {
	var cfg infra.Redis

	if err = typgo.ProcessConfig("REDIS", &cfg); err != nil {
		return
	}

	// TODO: using docker -it

	cmd := execkit.Command{
		Name: "redis-cli",
		Args: []string{
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-a", cfg.Password,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cmd.Print(os.Stdout)

	return cmd.Run(c.Ctx())
}

//
// redisDocker
//

var _ (typdocker.Composer) = (*redisDocker)(nil)

func (*redisDocker) DockerCompose() *typdocker.Recipe {
	var cfg infra.Redis
	if err := typgo.ProcessConfig("REDIS", &cfg); err != nil {
		panic("redis-docker: " + err.Error())
	}

	redis := dockerrx.Redis{
		Version:  typdocker.V3,
		Name:     "redis",
		Image:    "redis:4.0.5-alpine",
		Password: cfg.Password,
		Port:     cfg.Port,
	}

	return redis.DockerCompose()

}
