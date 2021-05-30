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
	t.initDefault()
	return &typgo.Task{
		Name:  t.Name,
		Usage: "redis tool",
		SubTasks: []*typgo.Task{
			{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
		},
	}
}

func (t *RedisTool) initDefault() {
	if t.Name == "" {
		t.Name = "REDIS"
	}
	if t.EnvKeys == nil {
		t.EnvKeys = EnvKeysWithPrefix(t.Name)
	}
	if t.DockerName == "" {
		t.DockerName = typgo.ProjectName + "-" + t.Name
	}
}

func (t *RedisTool) Console(c *typgo.Context) error {
	cfg := t.config()
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.DockerName,
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

func (t *RedisTool) config() *Config {
	return t.EnvKeys.Config()
}
