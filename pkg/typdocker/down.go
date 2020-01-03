package typdocker

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func downCmd() *cli.Command {
	return &cli.Command{
		Name:  "down",
		Usage: "Take down all docker containers according docker-compose",
		Action: func(ctx *cli.Context) (err error) {
			cmd := exec.Command("docker-compose", "down")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			return cmd.Run()
		},
	}
}
