package infra

import (
	"time"
)

type (
	// App is application configuration
	App struct {
		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
		Debug   bool   `default:"true"`
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
		DBName   string `required:"true" default:"MyLibrary"`
		User     string `required:"true" default:"postgres"`
		Password string `required:"true" default:"pgpass"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
	}
)
