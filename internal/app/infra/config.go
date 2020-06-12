package infra

import (
	"time"

	"github.com/typical-go/typical-rest-server/pkg/typpg"
)

type (
	// App is application configuration
	App struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `default:"false"`
	}
	// Redis Configuration
	Redis struct {
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
	// Pg is postgres configuration
	Pg struct {
		DBName   string `required:"true"`
		User     string `required:"true" default:"postgres"`
		Password string `required:"true" default:"pgpass"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
	}
)

//
// Pg
//

var _ typpg.Config = (*Pg)(nil)

// GetDBName to get database name
func (c *Pg) GetDBName() string { return c.DBName }

// GetUser to get user
func (c *Pg) GetUser() string { return c.User }

// GetPassword to get password
func (c *Pg) GetPassword() string { return c.Password }

// GetHost to get host
func (c *Pg) GetHost() string { return c.Host }

// GetPort to get port
func (c *Pg) GetPort() string { return c.Port }
