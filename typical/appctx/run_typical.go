package appctx

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// RunTypical to start the command line interface
func (c *Context) RunTypical(arguments []string) error {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Version

	app.Commands = c.standardTypicalCommand()
	for key := range c.Modules {
		app.Commands = append(app.Commands, c.Modules[key].Command())
	}
	return app.Run(arguments)
}

func (c *Context) standardTypicalCommand() []cli.Command {
	return []cli.Command{
		{Name: "build", Usage: "Build the binary", Action: notImplemented},
		{Name: "test", Usage: "Run the Test", Action: notImplemented},
		{Name: "run", Usage: "Run the binary", Action: notImplemented},
		{Name: "release", Usage: "Release the distribution", Action: notImplemented},
		{Name: "mock", Usage: "Generate mock class", Action: notImplemented},
		{Name: "readme", Usage: "Generate readme", Action: c.generateReadme},
	}
}

// NotImplement for not implemented function
func notImplemented(ctx *cli.Context) {
	fmt.Println("Not implemented")
}
