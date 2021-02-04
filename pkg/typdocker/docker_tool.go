package typdocker

import (
	"fmt"
	"os"
	"strings"

	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

type (
	// DockerTool for docker
	DockerTool struct{}
)

var (
	// DockerComposeYml is yml file
	DockerComposeYml = "docker-compose.yml"
	// Version of docker compose
	Version = "3"
)

//
// Command
//

var _ typgo.Tasker = (*DockerTool)(nil)

// Task for docker
func (m *DockerTool) Task() *typgo.Task {
	return &typgo.Task{
		Name:  "docker",
		Usage: "Docker utility",
		SubTasks: []*typgo.Task{
			{
				Name:    "up",
				Aliases: []string{"start"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "wipe"},
				},
				Usage:  "Spin up docker containers according docker-compose",
				Action: typgo.NewAction(DockerUp),
			},
			{
				Name:    "down",
				Aliases: []string{"stop"},
				Usage:   "Take down all docker containers according docker-compose",
				Action:  typgo.NewAction(DockerDown),
			},
			{
				Name:   "wipe",
				Usage:  "Kill all running docker container",
				Action: typgo.NewAction(DockerWipe),
			},
		},
	}
}

// DockerWipe clean all docker process
func DockerWipe(c *typgo.Context) error {
	ids, err := dockerIDs(c)
	if err != nil {
		return fmt.Errorf("Docker-ID: %w", err)
	}
	for _, id := range ids {
		if err := kill(c, id); err != nil {
			return fmt.Errorf("Fail to kill #%s: %s", id, err.Error())
		}
	}
	return nil
}

// DockerUp docker up
func DockerUp(c *typgo.Context) (err error) {
	if c.Bool("wipe") {
		if err := DockerWipe(c); err != nil {
			return err
		}
	}
	return c.Execute(&typgo.Bash{
		Name:   "docker-compose",
		Args:   []string{"up", "--remove-orphans", "-d"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

// DockerDown docker down
func DockerDown(c *typgo.Context) error {
	return c.Execute(&typgo.Bash{
		Name:   "docker-compose",
		Args:   []string{"down", "-v"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func dockerIDs(c *typgo.Context) (ids []string, err error) {
	var out strings.Builder

	if err = c.Execute(&typgo.Bash{
		Name:   "docker",
		Args:   []string{"ps", "-q"},
		Stderr: os.Stderr,
		Stdout: &out,
	}); err != nil {
		return
	}

	for _, id := range strings.Split(out.String(), "\n") {
		if id != "" {
			ids = append(ids, id)
		}
	}
	return
}

func kill(c *typgo.Context, id string) (err error) {
	return c.Execute(&typgo.Bash{
		Name:   "docker",
		Args:   []string{"kill", id},
		Stderr: os.Stderr,
	})
}
