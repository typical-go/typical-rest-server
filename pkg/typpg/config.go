package typpg

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

// Conn return connection string
func Conn(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// AdminConn return connection string for admin
func AdminConn(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/template1?sslmode=disable",
		c.User, c.Password, c.Host, c.Port)
}

// Configuration of postgres
func Configuration(s *Settings) *typgo.Configuration {
	if s == nil {
		panic("pg: configuration missing settings")
	}
	return &typgo.Configuration{
		Ctor: s.Ctor,
		Name: s.ConfigName,
		Spec: &Config{
			DBName:   s.DBName,
			User:     s.User,
			Password: s.Password,
			Host:     s.Host,
			Port:     s.Port,
		},
	}
}
