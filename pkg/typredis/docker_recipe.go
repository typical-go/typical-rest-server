package typredis

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

var (
	// DefaultDockerName is default value for docker name
	DefaultDockerName = "redis"

	// DefaultDockerImage is default value for docker image
	DefaultDockerImage = "redis:4.0.5-alpine"
)

// DockerRecipeV3 return recipe for docker-compose v3
func DockerRecipeV3() *typdocker.Recipe {
	return &typdocker.Recipe{
		Version: typdocker.V3,
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
