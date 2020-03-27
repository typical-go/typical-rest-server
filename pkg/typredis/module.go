package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-go/pkg/typapp"
)

// Module of Redis
type Module struct {
	dockerName  string
	dockerImage string
}

// New instance of redis module
func New() *Module {
	return &Module{
		dockerImage: "redis:4.0.5-alpine",
		dockerName:  "redis",
	}
}

// WithdockerImage to return module with new docker image
func (m *Module) WithdockerImage(dockerImage string) *Module {
	m.dockerImage = dockerImage
	return m
}

// WithdockerName to return module with new docker name
func (m *Module) WithdockerName(dockerName string) *Module {
	m.dockerName = dockerName
	return m
}

// Provide dependencies
func (m *Module) Provide() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(Connect),
	}
}

// Prepare the module
func (m *Module) Prepare() []*typapp.Preparation {
	return []*typapp.Preparation{
		typapp.NewPreparation(Ping),
	}
}

// Destroy dependencies
func (m *Module) Destroy() []*typapp.Destruction {
	return []*typapp.Destruction{
		typapp.NewDestruction(Disconnect),
	}
}

// Connect to redis server
func Connect(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:           cfg.Password,
		DB:                 cfg.DB,
		PoolSize:           cfg.PoolSize,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadWriteTimeout,
		WriteTimeout:       cfg.ReadWriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		MaxConnAge:         cfg.MaxConnAge,
	})
}

// Ping redis server
func Ping(client *redis.Client) (err error) {
	if err = client.Ping().Err(); err != nil {
		return fmt.Errorf("Redis: Ping: %w", err)
	}
	return
}

// Disconnect from service server
func Disconnect(client *redis.Client) (err error) {
	if err = client.Close(); err != nil {
		return fmt.Errorf("Redis: Disconnect: %w", err)
	}
	return
}
