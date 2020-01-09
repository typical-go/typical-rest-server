package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

// DockerCompose template
func (m *Module) DockerCompose(version typdocker.Version) *typdocker.ComposeObject {
	if version.IsV3() {
		return &typdocker.ComposeObject{
			Services: typdocker.Services{
				"postgres": typdocker.Service{
					Image: "postgres",
					Environment: map[string]string{
						"POSTGRES":          m.User,
						"POSTGRES_PASSWORD": m.Password,
						"PGDATA":            "/data/postgres",
					},
					Volumes:  []string{"postgres:/data/postgres"},
					Ports:    []string{fmt.Sprintf("%d:5432", m.Port)},
					Networks: []string{"postgres"},
					Restart:  "unless-stopped",
				},
			},
			Networks: typdocker.Networks{
				"postgres": typdocker.Network{
					Driver: "bridge",
				},
			},
			Volumes: typdocker.Volumes{
				"postgres": nil,
			},
		}
	}

	// TODO: create docker-compose for v2
	return nil
}
