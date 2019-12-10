package typredis

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-go/pkg/utility/envfile"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
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

// Module of Redis
type Module struct{}

// Configure Redis
func (r *Module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "REDIS"
	spec = &Config{}
	loadFn = func(loader typobj.Loader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (r *Module) Provide() []interface{} {
	return []interface{}{
		r.connect,
	}
}

// Prepare the module
func (r *Module) Prepare() []interface{} {
	return []interface{}{
		r.ping,
	}
}

// Destroy dependencies
func (r *Module) Destroy() []interface{} {
	return []interface{}{
		r.disconnect,
	}
}

// BuildCommands of module
func (r *Module) BuildCommands(c typobj.Cli) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return envfile.Load()
			},
			Subcommands: []*cli.Command{
				{Name: "console", Aliases: []string{"c"}, Usage: "Redis Interactive", Action: c.Action(r.console)},
			},
		},
	}
}

// DockerCompose template
func (r *Module) DockerCompose() typdocker.Compose {
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

func (*Module) connect(cfg Config) (client *redis.Client) {
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

func (*Module) ping(client *redis.Client) error {
	log.Info("Ping to Redis")
	return client.Ping().Err()
}

func (*Module) disconnect(client *redis.Client) (err error) {
	return client.Close()
}

func (*Module) console(config *Config) (err error) {
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
