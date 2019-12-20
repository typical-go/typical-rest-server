package typredis

import "github.com/typical-go/typical-rest-server/pkg/typdocker"

// DockerCompose template
func (*module) DockerCompose() typdocker.Compose {
	return typdocker.Compose{
		Services: map[string]interface{}{
			"redis": typdocker.Service{
				Image:    "redis:4.0.5-alpine",
				Command:  `redis-server --requirepass "${REDIS_PASSOWORD:-redispass}"`,
				Ports:    []string{`6379:6379`},
				Networks: []string{"redis"},
				Volumes:  []string{"redis:/data"},
			},
		},
		Networks: map[string]interface{}{
			"redis": nil,
		},
		Volumes: map[string]interface{}{
			"redis": nil,
		},
	}
}
