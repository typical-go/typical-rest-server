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
				"redis": typdocker.Service{
					Image:    "redis:4.0.5-alpine",
					Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, m.Password),
					Ports:    []string{fmt.Sprintf("%s:6379", m.Port)},
					Networks: []string{"redis"},
					Volumes:  []string{"redis:/data"},
				},
			},
			Networks: typdocker.Networks{
				"redis": nil,
			},
			Volumes: typdocker.Volumes{
				"redis": nil,
			},
		}
	}

	// TODO: docker-compose for v2
	return nil
}
