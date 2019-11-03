package typiobj

import "github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj/docker"

// DockerComposer responsible to compose docker
type DockerComposer interface {
	DockerCompose() docker.Compose
}
