package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerRecipe return recipe for docker-compose
func DockerRecipe() *typdocker.Recipe {
	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			DefaultDockerName: typdocker.Service{
				Image: DefaultDockerImage,
				Environment: map[string]string{
					"POSTGRES":          DefaultUser,
					"POSTGRES_PASSWORD": DefaultPassword,
					"PGDATA":            "/data/postgres",
				},
				Volumes:  []string{"postgres:/data/postgres"},
				Ports:    []string{fmt.Sprintf("%d:5432", DefaultPort)},
				Networks: []string{DefaultDockerName},
				Restart:  "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			DefaultDockerName: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			DefaultDockerName: nil,
		},
	}

}
