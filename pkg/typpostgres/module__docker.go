package typpostgres

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
					Image: m.dockerImage,
					Environment: map[string]string{
						"POSTGRES":          m.user,
						"POSTGRES_PASSWORD": m.password,
						"PGDATA":            "/data/postgres",
					},
					Volumes:  []string{"postgres:/data/postgres"},
					Ports:    []string{fmt.Sprintf("%d:5432", m.port)},
					Networks: []string{m.dockerName},
					Restart:  "unless-stopped",
				},
			},
			Networks: typdocker.Networks{
				m.dockerName: typdocker.Network{
					Driver: "bridge",
				},
			},
			Volumes: typdocker.Volumes{
				m.dockerName: nil,
			},
		}
	}

	// TODO: create docker-compose for v2
	return nil
}
