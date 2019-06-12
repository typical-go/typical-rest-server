package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/db"
	"github.com/typical-go/typical-rest-server/typical/projctx"
	"github.com/typical-go/typical-rest-server/typical/provider"
	"github.com/urfave/cli"
)

// Commands return list of command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "database",
			ShortName: "db",
			Subcommands: []cli.Command{
				{Name: "create", ShortName: "c", Usage: "Create New Database", Action: commandAction(db.Create)},
				{Name: "drop", ShortName: "d", Usage: "Drop Database", Action: commandAction(db.Drop)},
				{Name: "migrate", ShortName: "m", Usage: "Migrate Database", Action: commandAction(db.Migrate)},
				{Name: "rollback", ShortName: "r", Usage: "Rollback Database", Action: commandAction(db.Rollback)},
			},
		},

		{
			Name:      "config",
			ShortName: "cfg",
			Usage:     "Config details",
			Action: func(ctx *cli.Context) {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Type", "Required", "Default"})
				for _, detail := range config.Informations() {
					table.Append([]string{detail.Name, detail.Type, detail.Required, detail.Default})
				}
				table.Render()
			},
		},

		{
			Name:      "context",
			ShortName: "ctx",
			Action: func(ctx *cli.Context) {
				fmt.Println(projctx.String())
			},
		},

		{
			Name:      "dependency",
			ShortName: "dep",
			Subcommands: []cli.Command{
				{Name: "init", ShortName: "i", Usage: "Set up a new Go project, or migrate an existing one", Action: notImplement},
				{Name: "ensure", ShortName: "e", Usage: "install the project's dependencies", Action: notImplement},
				{Name: "update", ShortName: "u", Usage: "update the locked versions of all dependencies", Action: notImplement},
				{Name: "add", ShortName: "a", Usage: "add a dependency to the project", Action: notImplement},
			},
		},

		{
			Name:      "generate",
			ShortName: "gen",
			Subcommands: []cli.Command{
				{Name: "mock", Usage: "Generate mock", Action: notImplement},
				{Name: "migration", Usage: "Generate migration", Action: notImplement},
				{Name: "readme", Usage: "Generate readme", Action: notImplement},
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

func notImplement(ctx *cli.Context) {
	fmt.Println("Not implemented")
}
