package typicmd

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/yaml.v2"
)

// GenerateDockerCompose to generate docker-compose.yaml
func GenerateDockerCompose(ctx *typictx.ActionContext) (err error) {
	mainDocker := ctx.DockerCompose
	if mainDocker == nil {
		log.Info("No Docker Compose defined in Typical Context")
		return
	}
	for _, module := range ctx.Modules {
		moduleDocker := module.DockerCompose
		if moduleDocker == nil {
			continue
		}
		for _, name := range moduleDocker.ServiceKeys {
			mainDocker.RegisterService(name, moduleDocker.Services[name])
		}
		for _, name := range moduleDocker.NetworkKeys {
			mainDocker.RegisterNetwork(name, moduleDocker.Networks[name])
		}
		for _, name := range moduleDocker.VolumeKeys {
			mainDocker.RegisterVolume(name, moduleDocker.Volumes[name])
		}
	}
	d1, _ := yaml.Marshal(mainDocker)
	log.Info("Generate docker-compose.yml")
	return ioutil.WriteFile("docker-compose.yml", d1, 0644)
}

// DockerUp to create and start containers
func DockerUp(ctx *typictx.ActionContext) (err error) {
	if !ctx.Cli.Bool("no-compose") {
		err = GenerateDockerCompose(ctx)
		if err != nil {
			return
		}
	}
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// DockerDown to stop and remove containers, networks, images, and volumes
func DockerDown(ctx *typictx.ActionContext) (err error) {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
