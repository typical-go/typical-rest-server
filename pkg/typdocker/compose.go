package typdocker

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (m *Module) composeCmd(ctx *typcore.Context) *cli.Command {
	return &cli.Command{
		Name:  "compose",
		Usage: "Generate docker-compose.yaml",
		Action: func(c *cli.Context) (err error) {
			var (
				out []byte
				obj = dockerCompose(ctx.ProjectDescriptor, m.Version)
			)
			log.Info("Generate docker-compose.yml")
			if out, err = yaml.Marshal(obj); err != nil {
				return
			}
			if err = ioutil.WriteFile("docker-compose.yml", out, 0644); err != nil {
				return
			}
			return
		},
	}
}

func dockerCompose(d *typcore.ProjectDescriptor, version Version) (root *ComposeObject) {
	root = &ComposeObject{
		Version:  version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, module := range d.AllModule() {
		if composer, ok := module.(Composer); ok {
			if obj := composer.DockerCompose(version); obj != nil {
				root.Append(obj)
			}
		}
	}
	return
}
