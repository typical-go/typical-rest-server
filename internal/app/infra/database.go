package infra

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.uber.org/dig"

	// postgres driver
	_ "github.com/lib/pq"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type (
	dbConfigs struct {
		dig.In
		PgCfg    *DatabaseCfg `name:"pg"`
		MysqlCfg *DatabaseCfg `name:"mysql"`
	}
	// Databases setup output
	Databases struct {
		dig.Out
		Pg    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
	}
)

// NewDatabases return new instance of databases
// @ctor
func NewDatabases(c dbConfigs) Databases {
	return Databases{
		Pg:    createPGConn(c.PgCfg),
		MySQL: createMySQLConn(c.MysqlCfg),
	}
}

func createMySQLConn(p *DatabaseCfg) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&parseTime=true",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName))
	if err != nil {
		logrus.Fatalf("msyql: %s", err.Error())
	}
	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)
	if err = db.Ping(); err != nil {
		logrus.Fatalf("msyql: %s", err.Error())
	}
	return db
}

func createPGConn(p *DatabaseCfg) *sql.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		logrus.Fatalf("postgres: %s", err.Error())
	}

	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	if err = db.Ping(); err != nil {
		logrus.Fatalf("postgres: %s", err.Error())
	}

	return db
}
