package typdocker

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// New docker module
func New() *Module {
	return &Module{
		Version: "3",
	}
}

// Module of docker
type Module struct {
	Version   Version
	Composers []Composer
}

// WithVersion to set the version
func (m *Module) WithVersion(version Version) *Module {
	m.Version = version
	return m
}

// WithComposers to set the composers
func (m *Module) WithComposers(composers ...Composer) *Module {
	m.Composers = composers
	return m
}

// BuildCommands is command collection to called from
func (m *Module) BuildCommands(ctx *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Before: func(ctx *cli.Context) error {
				return typcfg.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				{
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
				},
				{
					Name:  "up",
					Usage: "Spin up docker containers according docker-compose",
					Action: func(ctx *cli.Context) (err error) {
						cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						return cmd.Run()
					},
				},
				{
					Name:  "down",
					Usage: "Take down all docker containers according docker-compose",
					Action: func(ctx *cli.Context) (err error) {
						cmd := exec.Command("docker-compose", "down")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						return cmd.Run()
					},
				},
				{
					Name:  "wipe",
					Usage: "Kill all running docker container",
					Action: func(ctx *cli.Context) (err error) {
						var builder strings.Builder
						cmd := exec.Command("docker", "ps", "-q")
						cmd.Stderr = os.Stderr
						cmd.Stdout = &builder
						if err = cmd.Run(); err != nil {
							return
						}
						dockerIDs := strings.Split(builder.String(), "\n")
						for _, id := range dockerIDs {
							if id != "" {
								cmd := exec.Command("docker", "kill", id)
								cmd.Stderr = os.Stderr
								if err = cmd.Run(); err != nil {
									return
								}
							}
						}
						return
					},
				},
			},
		},
	}
}
