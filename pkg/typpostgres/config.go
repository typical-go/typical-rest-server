package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typgo"
)

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// ConnStr return connection string
func (c *Config) ConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// Admin configuration
func (c *Config) Admin() *Config {
	return &Config{
		DBName:   "template1",
		User:     c.User,
		Password: c.Password,
		Host:     c.Host,
		Port:     c.Port,
	}
}

// Configuration of postgres
func Configuration(s *Setting) *typgo.Configuration {
	if s == nil {
		s = &Setting{}
	}
	return &typgo.Configuration{
		Name: GetConfigName(s),
		Spec: &Config{
			DBName:   GetDBName(s),
			User:     GetUser(s),
			Password: GetPassword(s),
			Host:     GetHost(s),
			Port:     GetPort(s),
		},
	}
}
