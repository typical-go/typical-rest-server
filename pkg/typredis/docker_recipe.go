package typredis

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// DockerRecipeV3 return recipe for docker-compose v3
func DockerRecipeV3(s *Settings) *typdocker.Recipe {
	return &typdocker.Recipe{
		Version: typdocker.V3,
		Services: typdocker.Services{
			s.DockerName: typdocker.Service{
				Image:    s.DockerName,
				Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, s.Password),
				Ports:    []string{fmt.Sprintf("%s:6379", s.Port)},
				Networks: []string{s.DockerName},
				Volumes:  []string{fmt.Sprintf("%s:/data", s.DockerName)},
			},
		},
		Networks: typdocker.Networks{
			s.DockerName: nil,
		},
		Volumes: typdocker.Volumes{
			s.DockerName: nil,
		},
	}

}
