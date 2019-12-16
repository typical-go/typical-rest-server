package typdocker

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func (d *docker) upCmd() *cli.Command {
	return &cli.Command{
		Name:   "up",
		Usage:  "Spin up docker containers according docker-compose",
		Action: d.up,
	}
}

func (d *docker) up(ctx *cli.Context) (err error) {
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
