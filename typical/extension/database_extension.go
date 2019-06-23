package extension

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// DatabaseExtension provide standard command for database operation like create, drop, migrate, rollback, generate migration
type DatabaseExtension struct {
	Extension
}

// Setup database extension
func (e *DatabaseExtension) Setup() error {
	return fmt.Errorf("not implemented")
}

//Command for database extension
func (e *DatabaseExtension) Command() cli.Command {
	return cli.Command{
		Name:      "database",
		ShortName: "db",
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: notImplement},
			{Name: "drop", Usage: "Drop Database", Action: notImplement},
			{Name: "migrate", Usage: "Migrate Database", Action: notImplement},
			{Name: "rollback", Usage: "Rollback Database", Action: notImplement},
		},
	}
}
