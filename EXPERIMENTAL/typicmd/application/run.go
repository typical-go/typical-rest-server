package application

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicli"

	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
)

// Run the application
func Run(ctx *typictx.Context) {
	ctxCli := &typicli.ContextCli{
		Context: ctx,
	}
	app := cli.NewApp()
	app.Name = ctx.Name
	app.Usage = ""
	app.Description = ctx.Description
	app.Version = ctx.Release.Version
	if runner, ok := ctx.AppModule.(typiobj.Runner); ok {
		app.Action = ctxCli.Action(runner.Run())
	}
	app.Before = typicli.LoadEnvFile
	if commander, ok := ctx.AppModule.(typicli.AppCommander); ok {
		app.Commands = commander.AppCommands(ctxCli)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
