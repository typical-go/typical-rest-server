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
			Usage:     "Run the server",
			Action: triggerAction(func(s *server) error {
				return s.Serve()
			}),
		},

		{
			Name:      "database",
			ShortName: "db",
			Usage:     "Database Administration",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "Create New Database",
					Action:    triggerAction(db.Create),
				},
				{
					Name:      "drop",
					ShortName: "d",
					Usage:     "Drop Database",
					Action:    triggerAction(db.Drop),
				},
				{
					Name:      "migrate",
					ShortName: "m",
					Usage:     "Migrate Database",
					Action:    triggerAction(db.Migrate),
				},
				{
					Name:      "rollback",
					ShortName: "r",
					Usage:     "Rollback Database",
					Action:    triggerAction(db.Rollback),
				},
			},
		},

		{
			Name:      "config",
			ShortName: "conf",
			Usage:     "Configuration",
			Action:    printConfigDetails,
		},

		// add more command here
	}
}

func triggerAction(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
