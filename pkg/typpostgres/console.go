package typpostgres

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdConsole(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:  "console",
		Usage: "PostgreSQL Interactive",
		Action: func(cliCtx *cli.Context) (err error) {
			return console(c.BuildContext(cliCtx))
		},
	}
}

func console(c *typbuildtool.BuildContext) (err error) {
	var cfg *Config
	if cfg, err = retrieveConfig(c); err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.CommandContext(c.Cli.Context, "psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
