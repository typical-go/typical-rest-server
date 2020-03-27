package typredis

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerCompose template
func (m *Module) DockerCompose(version typdocker.Version) *typdocker.ComposeObject {
	if version.IsV3() {
		return &typdocker.ComposeObject{
			Services: typdocker.Services{
				m.dockerName: typdocker.Service{
					Image:    m.dockerImage,
					Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, DefaultPassword),
					Ports:    []string{fmt.Sprintf("%s:6379", DefaultPort)},
					Networks: []string{m.dockerName},
					Volumes:  []string{fmt.Sprintf("%s:/data", m.dockerName)},
				},
			},
			Networks: typdocker.Networks{
				m.dockerName: nil,
			},
			Volumes: typdocker.Volumes{
				m.dockerName: nil,
			},
		}
	}

	// TODO: docker-compose for v2
	return nil
}
