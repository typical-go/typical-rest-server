package app

import (
	"os"

	"github.com/imantung/typical-go-server/config"
	"github.com/urfave/cli"
)

var (
	app  cli.App
	conf config.Config
)

// Run the service
func Run() (err error) {

	// prepare configuration
	conf, err = config.Load()
	if err != nil {
		return
	}

	// command line interface
	app := cli.NewApp()
	app.Name = config.App.Name
	app.Usage = config.App.Usage
	app.Version = config.App.Version

	initCommands(app)

	return app.Run(os.Args)
}
