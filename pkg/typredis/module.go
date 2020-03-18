package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typdep"
)

// Module of Redis
type Module struct {
	host        string
	port        string
	password    string
	dockerName  string
	dockerImage string
	prefix      string
}

// New instance of redis module
func New() *Module {
	return &Module{
		host:        "localhost",
		port:        "6379",
		password:    "redispass",
		dockerImage: "redis:4.0.5-alpine",
		dockerName:  "redis",
		prefix:      "REDIS",
	}
}

// WithHost to return module with new host
func (m *Module) WithHost(host string) *Module {
	m.host = host
	return m
}

// WithPort to return module with new port
func (m *Module) WithPort(port string) *Module {
	m.port = port
	return m
}

// Withpassword to return module with new password
func (m *Module) Withpassword(password string) *Module {
	m.password = password
	return m
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

// WithPrefix to return module with new prefix
func (m *Module) WithPrefix(prefix string) *Module {
	m.prefix = prefix
	return m
}

// Configure Redis
func (m *Module) Configure() *typcore.Configuration {
	return typcore.NewConfiguration(m.prefix, &Config{
		Host:     m.host,
		Port:     m.port,
		Password: m.password,
	})
}

// Provide dependencies
func (m *Module) Provide() []*typdep.Constructor {
	return []*typdep.Constructor{
		typdep.NewConstructor(m.connect),
	}
}

// Prepare the module
func (m *Module) Prepare() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(m.ping),
	}
}

// Destroy dependencies
func (m *Module) Destroy() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(m.disconnect),
	}
}

func (*Module) connect(cfg *Config) *redis.Client {
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

func (*Module) ping(client *redis.Client) (err error) {
	if err = client.Ping().Err(); err != nil {
		return fmt.Errorf("Redis: Ping: %w", err)
	}
	return
}

func (*Module) disconnect(client *redis.Client) (err error) {
	if err = client.Close(); err != nil {
		return fmt.Errorf("Redis: Disconnect: %w", err)
	}
	return
}
