package typpostgres

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const sourceURL = "file://scripts/db/migration"

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true" default:typical-rest`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

func openConnection(cfg *Config) (db *sql.DB, err error) {
	log.Info("Open postgres connection")
	db, err = sql.Open("postgres", dataSource(cfg))
	if err != nil {
		return
	}
	err = db.Ping()
	return
}

func closeConnection(db *sql.DB) error {
	log.Info("Close postgres connection")
	return db.Close()
}

func createDB(cfg *Config) (err error) {
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	conn, err := sql.Open("postgres", adminDataSource(cfg))
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}

func dropDB(cfg *Config) (err error) {
	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	conn, err := sql.Open("postgres", adminDataSource(cfg))
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}

func migrateDB(cfg *Config) error {
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	migration, err := migrate.New(sourceURL, dataSource(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(cfg *Config) error {
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	migration, err := migrate.New(sourceURL, dataSource(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func console(cfg *Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
