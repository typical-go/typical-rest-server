package typpg

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerRecipeV3 of postgres
func DockerRecipeV3(s *Settings) *typdocker.Recipe {
	if s == nil {
		panic("pg: docker-recipe missing settings")
	}

	name := s.DockerName
	image := s.DockerImage

	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			name: typdocker.Service{
				Image: image,
				Environment: map[string]string{
					"POSTGRES":          s.User,
					"POSTGRES_PASSWORD": s.Password,
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					"postgres:/data/postgres",
				},
				Ports: []string{
					fmt.Sprintf("%d:5432", s.Port),
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
