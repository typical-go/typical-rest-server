package typdocker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	// DockerTool is wrapper for docker-compose command support with predefined multiple compose file and environment
	DockerTool struct {
		ComposeFiles []string
		EnvFile      string
	}
)

//
// Command
//

var _ typgo.Tasker = (*DockerTool)(nil)

// Task for docker
func (m *DockerTool) Task() *typgo.Task {
	return &typgo.Task{
		Name:  "docker",
		Usage: "docker-compose wrapper",
		SubTasks: []*typgo.Task{
			{
				Name:   "wipe",
				Usage:  "Kill all running docker container",
				Action: typgo.NewAction(m.DockerWipe),
			},
		},
		Action: typgo.NewAction(m.DockerCompose),
	}
}

// DockerWipe clean all docker process
func (m *DockerTool) DockerWipe(c *typgo.Context) error {
	ids, err := m.dockerIDs(c)
	if err != nil {
		return fmt.Errorf("DockerTool-ID: %w", err)
	}
	for _, id := range ids {
		if err := m.kill(c, id); err != nil {
			return fmt.Errorf("fail to kill #%s: %s", id, err.Error())
		}
	}
	return nil
}

// DockerToolUp docker up
func (m *DockerTool) DockerCompose(c *typgo.Context) (err error) {
	var args []string
	args = append(args, "-p", typgo.ProjectName)
	if m.EnvFile != "" {
		args = append(args, "--env-file", m.EnvFile)
	}
	for _, file := range m.ComposeFiles {
		args = append(args, "-f", file)
	}
	args = append(args, c.Context.Args().Slice()...)

	return c.Execute(&typgo.Bash{
		Name:   "docker-compose",
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func (m *DockerTool) dockerIDs(c *typgo.Context) ([]string, error) {
	var out strings.Builder

	err := c.Execute(&typgo.Bash{
		Name:   "docker",
		Args:   []string{"ps", "-q"},
		Stderr: os.Stderr,
		Stdout: &out,
	})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range strings.Split(out.String(), "\n") {
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (m *DockerTool) kill(c *typgo.Context, id string) (err error) {
	return c.Execute(&typgo.Bash{
		Name:   "docker",
		Args:   []string{"kill", id},
		Stderr: os.Stderr,
	})
}

func ComposeFiles(dir string) []string {
	fileInfos, _ := ioutil.ReadDir(dir)
	var paths []string
	for _, info := range fileInfos {
		path := fmt.Sprintf("%s/%s", dir, info.Name())
		paths = append(paths, path)
	}
	return paths
}
