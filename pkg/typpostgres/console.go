package typpostgres

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m Module) consoleCmd(c typcore.Cli) *cli.Command {
	return &cli.Command{
		Name:   "console",
		Usage:  "PostgreSQL Interactive",
		Action: c.Action(m.console),
	}
}

func (Module) console(cfg Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
