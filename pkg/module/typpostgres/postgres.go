package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

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
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	migration, err := migrate.New(sourceURL, dataSource(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(cfg *Config) error {
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	migration, err := migrate.New(sourceURL, dataSource(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func seedDB(cfg *Config) (err error) {
	conn, err := sql.Open("postgres", dataSource(cfg))
	if err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(seedSrc)
	for _, f := range files {
		sqlFile := seedSrc + "/" + f.Name()
		log.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		b, err = ioutil.ReadFile(sqlFile)
		if err != nil {
			return
		}
		_, err = conn.Exec(string(b))
		if err != nil {
			return
		}
	}
	return
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
