package typredis

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	RedisTool struct {
		Name       string
		EnvKeys    *EnvKeys
		DockerName string
	}
)

var _ (typgo.Tasker) = (*RedisTool)(nil)

func (t *RedisTool) Task() *typgo.Task {
	return &typgo.Task{
		Name:  t.Name,
		Usage: "redis tool",
		SubTasks: []*typgo.Task{
			{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
		},
	}
}

func (t *RedisTool) Console(c *typgo.Context) error {
	cfg := t.config()
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.dockerName(),
			"redis-cli",
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-a", cfg.Pass,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

func (t *RedisTool) dockerName() string {
	dockerName := t.DockerName
	if dockerName == "" {
		dockerName = typgo.ProjectName + "-" + t.Name
	}
	return dockerName
}

func (t *RedisTool) config() *Config {
	return t.EnvKeys.Config()
}
