package app

import (
	"github.com/labstack/echo"
	"github.com/urfave/cli"
)

func initCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{Name: "Serve", ShortName: "s", Usage: "Serve the clients", Action: cmdServe},
		// add more command here
	}
}

func cmdServe(c *cli.Context) error {
	server := echo.New()

	initMiddlewares(server)
	initRoutes(server)

	return server.Start(conf.Address)
}
