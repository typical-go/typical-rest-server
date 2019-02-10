package app

import (
	"github.com/imantung/typical-go-server/db"
	"github.com/urfave/cli"
)

// Commands return list of command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Serve the clients",
			Action:    triggerAction(serve),
		},

		{
			Name:      "database",
			ShortName: "db",
			Usage:     "database administration",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "create new database",
					Action:    triggerAction(db.Create),
				},
				{
					Name:      "drop",
					ShortName: "d",
					Usage:     "drop database",
					Action:    triggerAction(db.Drop),
				},
				{
					Name:      "migrate",
					ShortName: "m",
					Usage:     "migrate database",
					Action:    triggerAction(db.Migrate),
				},
				{
					Name:      "rollback",
					ShortName: "r",
					Usage:     "rollback database",
					Action:    triggerAction(db.Rollback),
				},
			},
		},

		// add more command here
	}
}
