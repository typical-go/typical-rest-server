package main

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/_typical/project"
	"github.com/typical-go/typical-rest-server/db"
	"github.com/typical-go/typical-rest-server/provider"

	"github.com/urfave/cli"
)

// Commands return list of command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "database",
			ShortName: "db",
			Subcommands: []cli.Command{
				{Name: "create", ShortName: "c", Usage: "Create New Database", Action: invoke(db.Create)},
				{Name: "drop", ShortName: "d", Usage: "Drop Database", Action: invoke(db.Drop)},
				{Name: "migrate", ShortName: "m", Usage: "Migrate Database", Action: invoke(db.Migrate)},
				{Name: "rollback", ShortName: "r", Usage: "Rollback Database", Action: invoke(db.Rollback)},
			},
		},

		{
			Name:      "project",
			ShortName: "proj",
			Subcommands: []cli.Command{
				{Name: "config", ShortName: "cfg", Usage: "Config details", Action: print(project.ConfigDetail)},
				{Name: "context", ShortName: "ctx", Usage: "Context details", Action: print(project.ContextDetail)},
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

func invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := provider.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}

func notImplement(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func print(f func() string) interface{} {
	return func(ctx *cli.Context) error {
		fmt.Println(f())
		return nil
	}
}