package typdocker

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (m *Module) composeCmd(ctx *typcore.BuildContext) *cli.Command {
	return &cli.Command{
		Name:  "compose",
		Usage: "Generate docker-compose.yaml",
		Action: func(c *cli.Context) (err error) {
			if len(m.Composers) < 1 {
				return errors.New("No composers is set")
			}
			var out []byte
			log.Info("Generate docker-compose.yml")
			if out, err = yaml.Marshal(m.dockerCompose()); err != nil {
				return
			}
			if err = ioutil.WriteFile("docker-compose.yml", out, 0644); err != nil {
				return
			}
			return
		},
	}
}

func (m *Module) dockerCompose() (root *ComposeObject) {
	root = &ComposeObject{
		Version:  m.Version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, composer := range m.Composers {
		if obj := composer.DockerCompose(m.Version); obj != nil {
			root.Append(obj)
		}
	}
	return
}
