package typdocker

import (
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func (d *docker) wipeCmd() *cli.Command {
	return &cli.Command{
		Name:   "wipe",
		Usage:  "Kill all running docker container",
		Action: d.wipe,
	}
}

func (d *docker) wipe(ctx *cli.Context) (err error) {
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
}

type slice []string

func (w *slice) Write(data []byte) (n int, err error) {
	*w = append(*w, strings.TrimSpace(string(data)))
	return
}
