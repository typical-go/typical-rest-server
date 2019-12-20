package typredis

import (
	"os"
	"os/exec"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m *module) consoleCmd(c *typcore.Context) *cli.Command {
	return &cli.Command{
		Name:    "console",
		Aliases: []string{"c"},
		Usage:   "Redis Interactive",
		Action:  c.Action(m, m.console),
	}
}

func (*module) console(config *Config) (err error) {
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
