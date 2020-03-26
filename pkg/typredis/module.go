package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

const (
	// DefaultConfigName is default value for config name
	DefaultConfigName = "REDIS"
)

// Module of Redis
type Module struct {
	host        string
	port        string
	password    string
	dockerName  string
	dockerImage string
	configName  string
}

// New instance of redis module
func New() *Module {
	return &Module{
		host:        "localhost",
		port:        "6379",
		password:    "redispass",
		dockerImage: "redis:4.0.5-alpine",
		dockerName:  "redis",
		configName:  DefaultConfigName,
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

// WithConfigName to return module with new config name
func (m *Module) WithConfigName(configName string) *Module {
	m.configName = configName
	return m
}

// Configure Redis
func (m *Module) Configure() *typcfg.Configuration {
	return typcfg.NewConfiguration(m.configName, &Config{
		Host:     m.host,
		Port:     m.port,
		Password: m.password,
	})
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
