package app

import (
	"github.com/imantung/typical-go-server/app/server"
	"github.com/urfave/cli"
)

func initCommands(app *App) {
	app.Commands = []cli.Command{
		{
			Name:      "Serve",
			ShortName: "s",
			Usage:     "Serve the clients",
			Action:    cmdServe,
		},
	}
}

func cmdServe(c *cli.Context) (err error) {
	server := server.New()

	// TODO: get address from config
	return server.Start(":1323")
}
