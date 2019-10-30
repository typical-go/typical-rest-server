package application

import (
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
)

// Run the application
func Run(ctx *typictx.Context) {
	application := application{
		Context: ctx,
	}
	app := cli.NewApp()
	app.Name = ctx.Name
	app.Usage = ""
	app.Description = ctx.Description
	app.Version = ctx.Version
	app.Action = application.Run
	app.Before = func(ctx *cli.Context) error {
		return typienv.LoadEnvFile()
	}
	// for _, cmd := range c.Application.Commands {
	// 	cmd.Action = action(c, cmd.Action)
	// 	app.Commands = append(app.Commands, cmd)
	// }
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
