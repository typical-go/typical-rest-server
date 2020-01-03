package typdocker

import (
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func wipeCmd() *cli.Command {
	return &cli.Command{
		Name:  "wipe",
		Usage: "Kill all running docker container",
		Action: func(ctx *cli.Context) (err error) {
			var builder strings.Builder
			cmd := exec.Command("docker", "ps", "-q")
			cmd.Stderr = os.Stderr
			cmd.Stdout = &builder
			if err = cmd.Run(); err != nil {
				return
			}
			dockerIDs := strings.Split(builder.String(), "\n")
			for _, id := range dockerIDs {
				if id != "" {
					cmd := exec.Command("docker", "kill", id)
					cmd.Stderr = os.Stderr
					if err = cmd.Run(); err != nil {
						return
					}
				}
			}
			return
		},
	}
}
