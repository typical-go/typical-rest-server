package typredis

import (
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/typcfg"
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

	cmd := exec.CommandContext(c.Cli.Context, "redis-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func retrieveConfig() (*Config, error) {
	var cfg Config
	if err := typcfg.Process(DefaultConfigName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
