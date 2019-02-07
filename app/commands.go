package app

import (
	"github.com/imantung/typical-go-server/db"
	"github.com/urfave/cli"
)

func initCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Serve the clients",
			Action:    serve,
		},

		{
			Name:      "database",
			ShortName: "db",
			Usage:     "database comamnd",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "create new database",
					Action:    db.Create,
				},
				{
					Name:      "drop",
					ShortName: "d",
					Usage:     "drop database",
					Action:    db.Drop,
				},
				{
					Name:      "migrate",
					ShortName: "m",
					Usage:     "migrate database",
					Action:    db.Migrate,
				},
				{
					Name:      "rollback",
					ShortName: "r",
					Usage:     "rollback database",
					Action:    db.Rollback,
				},
			},
		},

		// add more command here
	}
}
