package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	// postgres driver
	_ "github.com/lib/pq"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
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
	// RedisCfg redis onfiguration
	// @envconfig (prefix:"REDIS")
	RedisCfg struct {
		Host     string `envconfig:"HOST" required:"true" default:"localhost"`
		Port     string `envconfig:"PORT" required:"true" default:"6379"`
		Password string `envconfig:"PASSWORD" default:"redispass"`
	}
	// PostgresCfg postgres configuration
	// @envconfig (prefix:"PG")
	PostgresCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"MyLibrary"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"pguser"`
		DBPass string `envconfig:"DBPASS" default:"pgpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"5432"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"30" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"6" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"30m" required:"true"`
	}
	// MySQLCfg is MySQL configuration
	// @envconfig (prefix:"MYSQL")
	MySQLCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"myalbum"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"mysql"`
		DBPass string `envconfig:"DBPASS" required:"true" default:"mypass"`
		Host   string `envconfig:"HOST" default:"localhost"`
		Port   string `envconfig:"PORT" default:"3306"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"30" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"6" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"30m" required:"true"`
	}
)

//
// RedisCfg
//

func (r *RedisCfg) createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
	})

	if err := client.Ping().Err(); err != nil {
		log.Fatalf("redis: %s", err.Error())
	}

	return client
}

//
// PostgresCfg
//

func (p *PostgresCfg) createConn() *sql.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("postgres: %s", err.Error())
	}

	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	if err = db.Ping(); err != nil {
		log.Fatalf("postgres: %s", err.Error())
	}

	return db
}

//
// MySQL
//

func (p *MySQLCfg) createConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&parseTime=true",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName))
	if err != nil {
		return nil, fmt.Errorf("mysql: %w", err)
	}
	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql: %w", err)
	}
	return db, nil
}
