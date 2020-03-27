package typredis

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerRecipe return recipe for docker-compose
func DockerRecipe() *typdocker.Recipe {

	return &typdocker.Recipe{
		Services: typdocker.Services{
			DefaultDockerName: typdocker.Service{
				Image:    DefaultDockerImage,
				Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, DefaultPassword),
				Ports:    []string{fmt.Sprintf("%s:6379", DefaultPort)},
				Networks: []string{DefaultDockerName},
				Volumes:  []string{fmt.Sprintf("%s:/data", DefaultDockerName)},
			},
		},
		Networks: typdocker.Networks{
			DefaultDockerName: nil,
		},
		Volumes: typdocker.Volumes{
			DefaultDockerName: nil,
		},
	}

}
