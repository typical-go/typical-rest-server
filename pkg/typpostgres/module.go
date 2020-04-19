package typpostgres

import (
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-go/pkg/typcfg"

	"github.com/typical-go/typical-go/pkg/typapp"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

// Module of postgres
func Module() *typapp.Module {
	return typapp.NewModule().
		Provide(
			typapp.NewConstructor(Connect),
		).
		Destroy(
			typapp.NewDestructor(Disconnect),
		).
		Prepare(typapp.NewPreparation(Ping)).
		Configure(typcfg.NewConfiguration(DefaultConfigName, DefaultConfig))
}

// Connect to postgres server
func Connect(cfg *Config) (db *sql.DB, err error) {
	if db, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		err = fmt.Errorf("Posgres: Connect: %w", err)
	}
	return
}

// Disconnect to postgres server
func Disconnect(db *sql.DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("Postgres: Disconnect: %w", err)
	}
	return
}

// Ping to postgres server
func Ping(db *sql.DB) (err error) {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Postgres: Ping: %w", err)
	}
	return
}

func dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
