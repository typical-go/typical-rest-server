package infra

import (
	"time"
)

type (
	// AppCfg application configuration
	// @envconfig (prefix:"APP")
	AppCfg struct {
		Address      string        `envconfig:"ADDRESS" default:":8089" required:"true"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		Debug        bool          `envconfig:"DEBUG" default:"true"`
	}
	// CacheCfg cache onfiguration
	// @envconfig (prefix:"CACHE")
	CacheCfg struct {
		DefaultMaxAge time.Duration `envconfig:"DEFAULT_MAX_AGE" default:"30s"`
		PrefixKey     string        `envconfig:"PREFIX_KEY" default:"cache_"`
		RedisHost     string        `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
		RedisPort     string        `envconfig:"REDIS_PORT" required:"true" default:"6379"`
		RedisPass     string        `envconfig:"REDIS_PASS" default:"redispass"`
	}
	// DatabaseCfg is MySQL configuration
	// @envconfig (prefix:"MYSQL" ctor:"mysql")
	// @envconfig (prefix:"PG" ctor:"pg")
	DatabaseCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"dbname"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"dbuser"`
		DBPass string `envconfig:"DBPASS" required:"true" default:"dbpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"9999"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"30" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"6" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"30m" required:"true"`
	}
)
