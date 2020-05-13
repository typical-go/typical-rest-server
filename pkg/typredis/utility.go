package typredis

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

// Utility of redis
func Utility() typgo.Utility {
	return typgo.NewUtility(Commands)
}

// Commands of redis utility
func Commands(c *typgo.BuildTool) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis utility",
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Action:  c.ActionFunc("REDIS", console),
				},
			},
		},
	}
}

func console(c *typgo.Context) (err error) {
	var config *Config
	if config, err = retrieveConfig(); err != nil {
		return
	}

	args := []string{
		"-h", config.Host,
		"-p", config.Port,
	}
	if config.Password != "" {
		args = append(args, "-a", config.Password)
	}
	// TODO: using docker -it

	cmd := execkit.Command{
		Name:   "redis-cli",
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cmd.Print(os.Stdout)

	return cmd.Run(c.Cli.Context)
}

func retrieveConfig() (*Config, error) {
	var cfg Config
	if err := typgo.ProcessConfig(DefaultConfigName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
