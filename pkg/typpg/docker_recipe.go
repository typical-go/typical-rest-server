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

	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			s.DockerName: typdocker.Service{
				Image: s.DockerImage,
				Environment: map[string]string{
					"POSTGRES":          s.User,
					"POSTGRES_PASSWORD": s.Password,
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					fmt.Sprintf("%s:/data/postgres", s.DockerName),
				},
				Ports: []string{
					fmt.Sprintf("%d:5432", s.Port),
				},
				Networks: []string{
					s.DockerName,
				},
				Restart: "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			s.DockerName: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			s.DockerName: nil,
		},
	}
}
