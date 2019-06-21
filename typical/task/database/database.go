package database

import (
	"github.com/typical-go/typical-go/appx"
	"gopkg.in/urfave/cli.v1"

	// load file source driver
	_ "github.com/golang-migrate/migrate/source/file"
)

// Create database
func Create(dbInfra appx.DBInfra) error {
	return dbInfra.Create()
}

// Drop database
func Drop(dbInfra appx.DBInfra) error {
	return dbInfra.Drop()
}

// Migrate database
func Migrate(dbInfra appx.DBInfra, args cli.Args) error {
	return dbInfra.Migrate(args)
}

// Rollback database
func Rollback(dbInfra appx.DBInfra, args cli.Args) error {
	return dbInfra.Rollback(args)
}
