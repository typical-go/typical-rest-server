package typtool

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	Redis struct {
		Name       string
		EnvKeys    *RedisEnvKeys
		DockerName string
	}
)

var _ (typgo.Tasker) = (*Redis)(nil)

func (t *Redis) Task() *typgo.Task {
	return &typgo.Task{
		Name:  t.Name,
		Usage: "redis tool",
		SubTasks: []*typgo.Task{
			{Name: "console", Usage: "Postgres console", Action: typgo.NewAction(t.Console)},
		},
	}
}

func (t *Redis) Console(c *typgo.Context) error {
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

func (t *Redis) dockerName() string {
	dockerName := t.DockerName
	if dockerName == "" {
		dockerName = typgo.ProjectName + "-" + t.Name
	}
	return dockerName
}

func (t *Redis) config() *RedisConfig {
	return t.EnvKeys.Config()
}
