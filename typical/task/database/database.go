package database

import (
	"github.com/typical-go/typical-rest-server/infra"
	"github.com/urfave/cli"

	// load file source driver
	_ "github.com/golang-migrate/migrate/source/file"
)

// Create database
func Create(dbInfra infra.DBInfra) error {
	return dbInfra.Create()
}

// Drop database
func Drop(dbInfra infra.DBInfra) error {
	return dbInfra.Drop()
}

// Migrate database
func Migrate(dbInfra infra.DBInfra, args cli.Args) error {
	return dbInfra.Migrate(args)
}

// Rollback database
func Rollback(dbInfra infra.DBInfra, args cli.Args) error {
	return dbInfra.Rollback(args)
}
