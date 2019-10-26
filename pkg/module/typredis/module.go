package typredis

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// Module of redis
func Module() *typictx.Module {
	return &typictx.Module{
		Name:     "Redis",
		Config:   typictx.NewConfig("REDIS", &Config{}),
		OpenFunc: Connect,
		Command: &typictx.Command{
			Name:       "redis",
			Usage:      "Redis Utility Tool",
			BeforeFunc: typienv.LoadEnvFile,
			SubCommands: []*typictx.Command{
				{Name: "console", ShortName: "c", ActionFunc: typictx.ActionFunction(Console)},
			},
		},
		DockerCompose: docker.NewCompose("").
			RegisterService("redis", &docker.Service{
				Image:   "redis:4.0.5-alpine",
				Command: `redis-server --requirepass "${REDIS_PASSOWORD:-redispass}"`,
				Ports: []string{
					`6379:6379`,
				},
				Networks: []string{
					"redis",
				},
				Volumes: []string{
					"redis:/data",
				},
			}).
			RegisterNetwork("redis", nil).
			RegisterVolume("redis", nil),
	}
}
