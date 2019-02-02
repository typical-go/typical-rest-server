package app

import (
	"github.com/urfave/cli"
)

// App application
type App struct {
	cli.App
}

// New create new app
func New() *App {
	app := &App{
		App: *cli.NewApp(),
	}
	app.Name = AppName
	app.Usage = AppUsage

	initCommands(app)

	return app
}
