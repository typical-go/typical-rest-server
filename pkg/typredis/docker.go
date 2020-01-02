package typredis

import "github.com/typical-go/typical-rest-server/pkg/typdocker"

// DockerCompose template
func (*Module) DockerCompose(version typdocker.Version) *typdocker.ComposeObject {
	if version.IsV3() {
		return &typdocker.ComposeObject{
			Services: typdocker.Services{
				"redis": typdocker.Service{
					Image:    "redis:4.0.5-alpine",
					Command:  `redis-server --requirepass "${REDIS_PASSOWORD:-redispass}"`,
					Ports:    []string{`6379:6379`},
					Networks: []string{"redis"},
					Volumes:  []string{"redis:/data"},
				},
			},
			Networks: typdocker.Networks{
				"redis": nil,
			},
			Volumes: typdocker.Volumes{
				"redis": nil,
			},
		}
	}

	// TODO: docker-compose for v2
	return nil
}
