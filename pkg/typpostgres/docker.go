package typpostgres

import "github.com/typical-go/typical-rest-server/pkg/typdocker"

// DockerCompose template
func (m *Module) DockerCompose() typdocker.Compose {
	return typdocker.Compose{
		Services: map[string]interface{}{
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
		Networks: map[string]interface{}{
			"postgres": typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: map[string]interface{}{
			"postgres": nil,
		},
	}
}
