package typredis

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli"
)

// Config is Redis Configuration
type Config struct {
	Host     string `required:"true" default:"localhost"`
	Port     string `required:"true" default:"6379"`
	Password string `default:"redispass"`
	DB       int    `default:"0"`

	PoolSize           int           `envconfig:"POOL_SIZE"  default:"20" required:"true"`
	DialTimeout        time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s" required:"true"`
	ReadWriteTimeout   time.Duration `envconfig:"READ_WRITE_TIMEOUT" default:"3s" required:"true"`
	IdleTimeout        time.Duration `envconfig:"IDLE_TIMEOUT" default:"5m" required:"true"`
	IdleCheckFrequency time.Duration `envconfig:"IDLE_CHECK_FREQUENCY" default:"1m" required:"true"`
	MaxConnAge         time.Duration `envconfig:"MAX_CONN_AGE" default:"30m" required:"true"`
}

// Module of redis
func Module() interface{} {
	return &redisModule{
		Name: "Redis",
		Configuration: typcfg.Configuration{
			Prefix: "REDIS",
			Spec:   &Config{},
		},
	}
}

type redisModule struct {
	typcfg.Configuration
	Name string
}

// BuildCommand of module
func (r redisModule) Command(c typcli.Cli) cli.Command {
	return cli.Command{
		Name:   "redis",
		Usage:  "Redis Tool",
		Before: typcli.LoadEnvFile,
		Subcommands: []cli.Command{
			{Name: "console", ShortName: "c", Usage: "Redis Interactive", Action: c.Action(r.console)},
		},
	}
}

// Provide dependencies
func (r redisModule) Provide() []interface{} {
	return []interface{}{
		r.loadConfig,
		r.connect,
	}
}

// Prepare the module
func (r redisModule) Prepare() []interface{} {
	return []interface{}{
		r.ping,
	}
}

// Destroy dependencies
func (r redisModule) Destroy() []interface{} {
	return []interface{}{
		r.disconnect,
	}
}

func (r redisModule) DockerCompose() typdocker.Compose {
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

func (r redisModule) loadConfig(loader typcfg.Loader) (cfg Config, err error) {
	err = loader.Load(r.Configuration, &cfg)
	return
}

func (redisModule) connect(cfg Config) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
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
	return
}

func (redisModule) ping(client *redis.Client) error {
	log.Info("Ping to Redis")
	return client.Ping().Err()
}

func (redisModule) disconnect(client *redis.Client) (err error) {
	return client.Close()
}

func (redisModule) console(config *Config) (err error) {
	args := []string{
		"-h", config.Host,
		"-p", config.Port,
	}
	if config.Password != "" {
		args = append(args, "-a", config.Password)
	}
	// TODO: using docker -it
	cmd := exec.Command("redis-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
