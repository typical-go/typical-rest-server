package db

import (
	"database/sql"
	"fmt"

	"github.com/imantung/go-helper/dbkit"
	"github.com/imantung/typical-go-server/config"
)

// Create database
func Create(conf config.Config) error {
	query := fmt.Sprintf(`CREATE DATABASE %s`, conf.DbName)
	fmt.Println(query)
	return executeFromTemplateDB(conf, query)
}

// Drop database
func Drop(conf config.Config) error {
	query := fmt.Sprintf(`DROP DATABASE %s`, conf.DbName)
	fmt.Println(query)
	return executeFromTemplateDB(conf, query)
}

// Migrate database
func Migrate(conn *sql.DB) error {
	fmt.Println("Migrate Database")
	return conn.Ping()
}

// Rollback database
func Rollback(conn *sql.DB) error {
	fmt.Println("Rollback Database")
	return conn.Ping()
}

func executeFromTemplateDB(conf config.Config, query string) (err error) {
	pgConf := dbkit.PgConfig{
		Host:     conf.DbHost,
		Port:     conf.DbPort,
		DbName:   "template1",
		User:     conf.DbUser,
		Password: conf.DbPassword,
		SslMode:  "disable",
	}
	conn, err := sql.Open("postgres", pgConf.ConnectionString())
	if err != nil {
		return
	}
	_, err = conn.Exec(query)
	return
}
