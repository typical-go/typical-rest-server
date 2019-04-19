package app

import (
	"os"

	"github.com/imantung/typical-go-server/config"
	"github.com/imantung/typical-go-server/db"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// Commands return list of command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Run the server",
			Action: commandAction(func(s *server) error {
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
					Action:    commandAction(db.Create),
				},
				{
					Name:      "drop",
					ShortName: "d",
					Usage:     "Drop Database",
					Action:    commandAction(db.Drop),
				},
				{
					Name:      "migrate",
					ShortName: "m",
					Usage:     "Migrate Database",
					Action:    commandAction(db.Migrate),
				},
				{
					Name:      "rollback",
					ShortName: "r",
					Usage:     "Rollback Database",
					Action:    commandAction(db.Rollback),
				},
			},
		},

		{
			Name:      "config",
			ShortName: "conf",
			Usage:     "Configuration",
			Action: func(ctx *cli.Context) {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Type", "Required", "Default"})
				for _, detail := range config.Information() {
					table.Append([]string{detail.Name, detail.Type, detail.Required, detail.Default})
				}
				table.Render()
			},
		},

		// add more command here
	}
}

func commandAction(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
