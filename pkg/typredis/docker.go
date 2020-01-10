package typredis

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

// DockerCompose template
func (m *Module) DockerCompose(version typdocker.Version) *typdocker.ComposeObject {
	if version.IsV3() {
		return &typdocker.ComposeObject{
			Services: typdocker.Services{
				m.DockerName: typdocker.Service{
					Image:    m.DockerImage,
					Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, m.Password),
					Ports:    []string{fmt.Sprintf("%s:6379", m.Port)},
					Networks: []string{m.DockerName},
					Volumes:  []string{fmt.Sprintf("%s:/data", m.DockerName)},
				},
			},
			Networks: typdocker.Networks{
				m.DockerName: nil,
			},
			Volumes: typdocker.Volumes{
				m.DockerName: nil,
			},
		}
	}

	// TODO: docker-compose for v2
	return nil
}
