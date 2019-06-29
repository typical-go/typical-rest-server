package xdb

import (
	"fmt"

	"github.com/typical-go/typical-go/appx"
	"gopkg.in/urfave/cli.v1"

	// load file source driver
	_ "github.com/golang-migrate/migrate/source/file"
)

// DefaultMigrationDirectory default migration diremake ctory
var DefaultMigrationDirectory = "scripts/migration"

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
	source := migrationSource(args)
	return dbInfra.Migrate(source)
}

// Rollback database
func Rollback(dbInfra appx.DBInfra, args cli.Args) error {
	source := migrationSource(args)
	return dbInfra.Rollback(source)
}

func migrationSource(args cli.Args) string {
	dir := DefaultMigrationDirectory
	if len(args) > 0 {
		dir = args.First()
	}
	return fmt.Sprintf("file://%s", dir)
}
