package typpostgres

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typcfg"
)

func dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}

func retrieveConfig() (*Config, error) {
	var cfg Config
	if err := typcfg.Process(DefaultConfigName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
