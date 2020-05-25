package typredis

import (
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

type utility struct {
	*Settings
}

// Utility of redis
func Utility(s *Settings) typgo.Utility {
	return &utility{
		Settings: s,
	}
}

// Commands of redis utility
func (u *utility) Commands(c *typgo.BuildCli) []*cli.Command {
	name := u.UtilityCmd
	return []*cli.Command{
		{
			Name:  name,
			Usage: "Redis utility",
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Action:  c.ActionFn(name, u.console),
				},
			},
		},
	}
}

func (u *utility) console(c *typgo.Context) (err error) {
	var config *Config
	if config, err = u.retrieveConfig(); err != nil {
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

	return cmd.Run(c.Ctx())
}

func (u *utility) retrieveConfig() (*Config, error) {
	var cfg Config
	if err := typgo.ProcessConfig(u.ConfigName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
