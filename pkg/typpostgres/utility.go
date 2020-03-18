package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
)

func (m *Postgres) reset(c *Context) (err error) {
	if err = m.drop(c); err != nil {
		return
	}
	if err = m.create(c); err != nil {
		return
	}
	if err = m.migrate(c); err != nil {
		return
	}
	if err = m.seed(c); err != nil {
		return
	}
	return
}

func (m *Postgres) create(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = c.Config(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func (m *Postgres) drop(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = c.Config(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func (*Postgres) console(c *Context) (err error) {
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (m *Postgres) migrate(c *Context) (err error) {
	var (
		migration *migrate.Migrate
		sourceURL = "file://" + m.migrationSource
		cfg       *Config
	)

	if cfg, err = c.Config(); err != nil {
		return
	}

	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (m *Postgres) rollback(c *Context) (err error) {
	var migration *migrate.Migrate
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}

	sourceURL := "file://" + m.migrationSource
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (m *Postgres) seed(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}
	if conn, err = sql.Open("postgres", m.dataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(m.seedSource)
	for _, f := range files {
		sqlFile := m.seedSource + "/" + f.Name()
		log.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			log.Error(err.Error())
			continue
		}
		if _, err = conn.Exec(string(b)); err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return
}
