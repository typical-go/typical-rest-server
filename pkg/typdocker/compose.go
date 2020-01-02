package typdocker

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (d *docker) composeCmd() *cli.Command {
	return &cli.Command{
		Name:   "compose",
		Usage:  "Generate docker-compose.yaml",
		Action: d.compose,
	}
}

func (d *docker) compose(ctx *cli.Context) (err error) {
	log.Info("Generate docker-compose.yml")
	out, _ := yaml.Marshal(d.dockerCompose())
	return ioutil.WriteFile("docker-compose.yml", out, 0644)
}

func (d *docker) dockerCompose() (root *ComposeObject) {
	root = &ComposeObject{
		Version:  d.Version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, module := range d.AllModule() {
		if composer, ok := module.(Composer); ok {
			if obj := composer.DockerCompose(d.Version); obj != nil {
				root.Append(obj)
			}
		}
	}
	return
}
