package typdocker

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typctx"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type dockerCommand struct {
	*typctx.Context
}

func (c dockerCommand) Compose(ctx *cli.Context) (err error) {
	log.Info("Generate docker-compose.yml")
	out, _ := yaml.Marshal(c.dockerCompose())
	return ioutil.WriteFile("docker-compose.yml", out, 0644)
}

func (c dockerCommand) Up(ctx *cli.Context) (err error) {
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (c dockerCommand) Down(ctx *cli.Context) (err error) {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// DockerCompose get docker compose
func (c dockerCommand) dockerCompose() (dc Compose) {
	dc.Version = "3"
	dc.Services = make(map[string]interface{})
	dc.Networks = make(map[string]interface{})
	dc.Volumes = make(map[string]interface{})
	for _, module := range c.AllModule() {
		if composer, ok := module.(DockerComposer); ok {
			dc.Add(composer.DockerCompose())
		}
	}
	return
}
