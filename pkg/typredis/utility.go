package typredis

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

// Utility return new instance of PostgresUtility
func Utility() typbuildtool.Utility {
	return typbuildtool.NewUtility(Commands)
}

// Commands of redis utility
func Commands(c *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis utility",
			Subcommands: []*cli.Command{
				{
					Name: "console",
					Action: func(cliCtx *cli.Context) (err error) {
						return console(c)
					},
				},
			},
		},
	}
}

func console(c *typbuildtool.Context) (err error) {
	var config *Config
	if config, err = retrieveConfig(c); err != nil {
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
	cmd := exec.Command("redis-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func retrieveConfig(c *typbuildtool.Context) (cfg *Config, err error) {
	var v interface{}
	var ok bool

	if v, err = c.RetrieveConfig(DefaultConfigName); err != nil {
		return
	}

	if cfg, ok = v.(*Config); !ok {
		return nil, fmt.Errorf("Redis: Get config for '%s' but invalid type", DefaultConfigName)
	}

	return
}
