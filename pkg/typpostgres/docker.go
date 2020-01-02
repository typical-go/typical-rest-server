package typpostgres

import "github.com/typical-go/typical-rest-server/pkg/typdocker"

// DockerCompose template
func (m *Module) DockerCompose(version typdocker.Version) *typdocker.ComposeObject {
	if version.IsV3() {
		return &typdocker.ComposeObject{
			Services: typdocker.Services{
				"postgres": typdocker.Service{
					Image: "postgres",
					Environment: map[string]string{
						"POSTGRES":          "${PG_USER:-postgres}",
						"POSTGRES_PASSWORD": "${PG_PASSWORD:-pgpass}",
						"PGDATA":            "/data/postgres",
					},
					Volumes:  []string{"postgres:/data/postgres"},
					Ports:    []string{"${PG_PORT:-5432}:5432"},
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
