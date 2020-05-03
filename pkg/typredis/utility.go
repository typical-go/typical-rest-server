package typredis

import (
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

// Utility return new instance of PostgresUtility
func Utility() typbuildtool.Utility {
	return typbuildtool.NewUtility(Commands).
		Configure(&typcfg.Configuration{Name: DefaultConfigName, Spec: DefaultConfig})
}

// Commands of redis utility
func Commands(c *typbuildtool.Context) []*cli.Command {
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

func console(c *typbuildtool.CliContext) (err error) {
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
	cmd := exec.CommandContext(c.Context, "redis-cli", args...)
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
