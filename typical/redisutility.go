package typical

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

type redisUtility struct{}

var _ typgo.Utility = (*redisUtility)(nil)

func (*redisUtility) Commands(c *typgo.BuildCli) ([]*cli.Command, error) {
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
	}, nil
}

func redisConsole(c *typgo.Context) error {
	return c.Execute(&execkit.Command{
		Name: "redis-cli",
		Args: []string{
			"-h", os.Getenv("REDIS_HOST"),
			"-p", os.Getenv("REDIS_PORT"),
			"-a", os.Getenv("REDIS_PASSWORD"),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})

}
