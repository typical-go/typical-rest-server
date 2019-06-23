package extension

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical/task/database"
	"gopkg.in/urfave/cli.v1"
)

// GoExtension provide standar go command like build, run, test, mock, etc
type GoExtension struct{}

// Setup go extension
func (e *GoExtension) Setup() error {
	return fmt.Errorf("not implemented")
}

//Command for go extension
func (e *GoExtension) Command() cli.Command {
	return cli.Command{
		Name: "go",
		Subcommands: []cli.Command{
			{Name: "build", Usage: "Build the binary", Action: invoke(database.Create)},
			{Name: "test", Usage: "Run the Test", Action: invoke(database.Drop)},
			{Name: "run", Usage: "Run the binary", Action: invoke(database.Migrate)},
			{Name: "distribute", Usage: "Create distribution package", Action: invoke(database.Migrate)},
			{Name: "mock", Usage: "Generate mock class", Action: invoke(database.Rollback)},
		},
	}
}
