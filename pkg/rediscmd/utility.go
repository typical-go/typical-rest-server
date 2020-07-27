package rediscmd

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

// Utility for redis
type Utility struct {
	Name        string
	HostEnv     string
	PortEnv     string
	PasswordEnv string
}

var _ typgo.Cmd = (*Utility)(nil)

func (u *Utility) validate() string {
	if u.Name == "" {
		return "missing Name"
	}
	if u.HostEnv == "" {
		return "missing HostEnv"
	}
	if u.PortEnv == "" {
		return "missing PortEnv"
	}
	if u.PasswordEnv == "" {
		return "missing PasswordEnv"
	}
	return ""
}

// Command for utility
func (u *Utility) Command(sys *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:  u.Name,
		Usage: "Redis utility",
		Subcommands: []*cli.Command{
			{
				Name:    "console",
				Aliases: []string{"c"},
				Action:  sys.ActionFn(typgo.NewAction(u.console)),
			},
		},
	}
}

func (u *Utility) console(c *typgo.Context) error {
	if errMsg := u.validate(); errMsg != "" {
		return fmt.Errorf("redis-cmd: %s", errMsg)
	}
	return c.Execute(&execkit.Command{
		Name: "redis-cli",
		Args: []string{
			"-h", os.Getenv(u.HostEnv),
			"-p", os.Getenv(u.PortEnv),
			"-a", os.Getenv(u.PasswordEnv),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}
