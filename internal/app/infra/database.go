package infra

import (
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-rest-server/internal/app/infra/log"

	// postgres driver
	_ "github.com/lib/pq"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

func createMySQLConn(p *DatabaseCfg) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&parseTime=true",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName))
	if err != nil {
		log.Fatalf("msyql: %s", err.Error())
	}
	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)
	if err = db.Ping(); err != nil {
		log.Fatalf("msyql: %s", err.Error())
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
		log.Fatalf("postgres: %s", err.Error())
	}

	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	if err = db.Ping(); err != nil {
		log.Fatalf("postgres: %s", err.Error())
	}

	return db
}
