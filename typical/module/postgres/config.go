package postgres

import (
	"fmt"
)

// Config is postgres configuration
type Config struct {
	// dbtool.Config

	DbName       string `required:"true"`
	User         string `required:"true"`
	Password     string `required:"true"`
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	MigrationSrc string `default:"scripts/migration"`
}

// DataSource return connection string
func (c *Config) DataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DbName)
}

// AdminDataSource return connection string for adminitration script
func (c *Config) AdminDataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}

// DatabaseName return the database name
func (c *Config) DatabaseName() string {
	return c.DbName
}

// DriverName return the driver name
func (c *Config) DriverName() string {
	return "postgres"
}

// MigrationSource return the migration source
func (c *Config) MigrationSource() string {
	return c.MigrationSrc
}
