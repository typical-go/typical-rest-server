package extension

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// GoExtension provide standar go command like build, run, test, mock, etc
type GoExtension struct {
	Extension
	ActionTrigger
}

// Setup go extension
func (e *GoExtension) Setup() error {
	return fmt.Errorf("not implemented")
}

//Command for go extension
func (e *GoExtension) Command() cli.Command {
	return cli.Command{
		Name: "go",
		Subcommands: []cli.Command{
			{Name: "build", Usage: "Build the binary", Action: e.NotImplement},
			{Name: "test", Usage: "Run the Test", Action: e.NotImplement},
			{Name: "run", Usage: "Run the binary", Action: e.NotImplement},
			{Name: "distribute", Usage: "Create distribution package", Action: e.NotImplement},
			{Name: "mock", Usage: "Generate mock class", Action: e.NotImplement},
		},
	}
}
