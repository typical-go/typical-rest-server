package typdocker

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type dockerCommand struct {
	*typictx.Context
}

func (c dockerCommand) Compose(ctx *cli.Context) (err error) {
	log.Info("Generate docker-compose.yml")
	out, _ := yaml.Marshal(c.DockerCompose())
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
