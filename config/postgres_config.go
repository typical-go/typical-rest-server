package config

import "fmt"

// PostgresConfig contain postgres database configuration
type PostgresConfig struct {
	DbName       string `required:"true" default:"typical-rest-server"`
	User         string `required:"true" default:"default_user"`
	Password     string `required:"true" default:"default_password"`
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	MigrationSrc string `default:"scripts/migration"`
}

// DataSource return connection string
func (c *PostgresConfig) DataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DbName)
}

// AdminDataSource return connection string for adminitration script
func (c *PostgresConfig) AdminDataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}

// DatabaseName return the database name
func (c *PostgresConfig) DatabaseName() string {
	return c.DbName
}

// DriverName return the driver name
func (c *PostgresConfig) DriverName() string {
	return "postgres"
}

// MigrationSource return the migration source
func (c *PostgresConfig) MigrationSource() string {
	return c.MigrationSrc
}
