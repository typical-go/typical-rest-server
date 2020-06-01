package typical

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/urfave/cli/v2"
)

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

func redisDocker() *typdocker.Recipe {
	name := "redis"
	image := "redis:4.0.5-alpine"

	var cfg infra.Redis
	typgo.ProcessConfig("REDIS", &cfg)

	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			name: typdocker.Service{
				Image:    image,
				Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, cfg.Password),
				Ports:    []string{fmt.Sprintf("%s:6379", cfg.Port)},
				Networks: []string{name},
				Volumes:  []string{fmt.Sprintf("%s:/data", name)},
			},
		},
		Networks: typdocker.Networks{
			name: nil,
		},
		Volumes: typdocker.Volumes{
			name: nil,
		},
	}
}
