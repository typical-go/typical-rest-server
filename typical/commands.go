package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/db"
	"github.com/typical-go/typical-rest-server/typical/provider"
	"github.com/urfave/cli"
)

// Commands return list of command
func Commands() []cli.Command {
	return []cli.Command{
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
				for _, detail := range config.Informations() {
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
		container := provider.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
