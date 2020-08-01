package typdocker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"
)

var (
	// DockerComposeYml is yml file
	DockerComposeYml = "docker-compose.yml"
	// Version of docker compose
	Version = "3"
)

type (
	// DockerCmd for docker
	DockerCmd struct {
		Composers []Composer
	}
)

//
// Command
//

var _ typgo.Cmd = (*DockerCmd)(nil)
var _ typgo.Action = (*DockerCmd)(nil)

// Command of docker
func (m *DockerCmd) Command(sys *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "Docker utility",
		Subcommands: []*cli.Command{
			m.CmdCompose(sys),
			m.CmdUp(sys),
			m.CmdDown(sys),
			m.CmdWipe(sys),
		},
	}
}

// CmdCompose command compose
func (m *DockerCmd) CmdCompose(c *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:   "compose",
		Usage:  "Generate docker-compose.yaml",
		Action: c.ActionFn(m),
	}
}

// Execute docker command
func (m *DockerCmd) Execute(c *typgo.Context) error {
	if len(m.Composers) < 1 {
		return errors.New("Nothing to compose")
	}

	root, err := compile(Version, m.Composers)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	out, err := yaml.Marshal(root)
	if err != nil {
		return err
	}

	fmt.Println("Generate docker-compose.yml")
	return ioutil.WriteFile(DockerComposeYml, out, 0777)
}

// Compile recipes to yaml
func compile(version string, composers []Composer) (*Recipe, error) {
	root := &Recipe{
		Version:  version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, composer := range composers {
		obj, err := composer.ComposeV3()
		if err != nil {
			return nil, err
		}
		for k, service := range obj.Services {
			root.Services[k] = service
		}
		for k, network := range obj.Networks {
			root.Networks[k] = network
		}
		for k, volume := range obj.Volumes {
			root.Volumes[k] = volume
		}
	}
	return root, nil
}

// CmdWipe command wipe
func (m *DockerCmd) CmdWipe(c *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:   "wipe",
		Usage:  "Kill all running docker container",
		Action: c.ActionFn(typgo.NewAction(m.dockerWipe)),
	}
}

func (m *DockerCmd) dockerWipe(c *typgo.Context) error {
	ids, err := dockerIDs(c)
	if err != nil {
		return fmt.Errorf("Docker-ID: %w", err)
	}
	for _, id := range ids {
		if err := kill(c, id); err != nil {
			return fmt.Errorf("Fail to kill #%s: %s", id, err.Error())
		}
	}
	return nil
}

// CmdUp command up
func (m *DockerCmd) CmdUp(c *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:    "up",
		Aliases: []string{"start"},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "wipe"},
		},
		Usage:  "Spin up docker containers according docker-compose",
		Action: c.ActionFn(typgo.NewAction(m.dockerUp)),
	}
}

func (m *DockerCmd) dockerUp(c *typgo.Context) (err error) {
	if c.Bool("wipe") {
		if err := m.dockerWipe(c); err != nil {
			return err
		}
	}
	return c.Execute(&execkit.Command{
		Name:   "docker-compose",
		Args:   []string{"up", "--remove-orphans", "-d"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

// CmdDown command down
func (m *DockerCmd) CmdDown(c *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:    "down",
		Aliases: []string{"stop"},
		Usage:   "Take down all docker containers according docker-compose",
		Action:  c.ActionFn(typgo.NewAction(dockerDown)),
	}
}

func dockerDown(c *typgo.Context) error {
	return c.Execute(&execkit.Command{
		Name:   "docker-compose",
		Args:   []string{"down", "-v"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func dockerIDs(c *typgo.Context) (ids []string, err error) {
	var out strings.Builder

	if err = c.Execute(&execkit.Command{
		Name:   "docker",
		Args:   []string{"ps", "-q"},
		Stderr: os.Stderr,
		Stdout: &out,
	}); err != nil {
		return
	}

	for _, id := range strings.Split(out.String(), "\n") {
		if id != "" {
			ids = append(ids, id)
		}
	}
	return
}

func kill(c *typgo.Context, id string) (err error) {
	return c.Execute(&execkit.Command{
		Name:   "docker",
		Args:   []string{"kill", id},
		Stderr: os.Stderr,
	})
}
