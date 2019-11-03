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
	app := cli.NewApp()
	app.Name = ctx.Name
	app.Usage = ""
	app.Description = ctx.Description
	app.Version = ctx.Release.Version
	if runner, ok := ctx.AppModule.(typiobj.Runner); ok {
		app.Action = typicli.Action(ctx, runner.Run())
	}
	app.Before = typicli.LoadEnvFile
	if appCli, ok := ctx.AppModule.(typictx.AppCLI); ok {
		app.Commands = appCli.AppCommands(ctx)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
