package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerRecipeV3 of postgres
func DockerRecipeV3(s *Setting) *typdocker.Recipe {
	if s == nil {
		s = &Setting{}
	}
	name := GetDockerName(s)
	image := GetDockerImage(s)

	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			name: typdocker.Service{
				Image: image,
				Environment: map[string]string{
					"POSTGRES":          GetUser(s),
					"POSTGRES_PASSWORD": GetPassword(s),
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					"postgres:/data/postgres",
				},
				Ports: []string{
					fmt.Sprintf("%d:5432", GetPort(s)),
				},
				Networks: []string{
					name,
				},
				Restart: "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			name: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			name: nil,
		},
	}
}
