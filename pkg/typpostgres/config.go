package typpostgres

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
)

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// Configuration of Postgres
func Configuration() *typcfg.Configuration {
	return typcfg.NewConfiguration(DefaultConfigName, &Config{
		DBName:   DefaultUser,
		User:     DefaultUser,
		Password: DefaultPassword,
		Host:     DefaultHost,
		Port:     DefaultPort,
	})
}
