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

func (m *Postgres) reset(cfg Config) (err error) {
	if err = m.drop(cfg); err != nil {
		return
	}
	if err = m.create(cfg); err != nil {
		return
	}
	if err = m.migrate(cfg); err != nil {
		return
	}
	if err = m.seed(cfg); err != nil {
		return
	}
	return
}

func (m *Postgres) create(cfg Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}
	_, err = conn.Exec(query)
	return
}

func (m *Postgres) drop(cfg Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}

func (*Postgres) console(cfg Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (m *Postgres) migrate(cfg Config) (err error) {
	var (
		migration *migrate.Migrate
		sourceURL = "file://" + m.migrationSource
	)
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (m *Postgres) rollback(cfg Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + m.migrationSource
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, m.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (m *Postgres) seed(cfg Config) (err error) {
	var conn *sql.DB
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
