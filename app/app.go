package app

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli"
)

var (
	app  cli.App
	conf config
)

// Run run the application
func Run() (err error) {

	// prepare configuration
	err = envconfig.Process(appConfigPrefix, &conf)
	if err != nil {
		return
	}

	// command line interface
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage

	initCommands(app)

	return app.Run(os.Args)
}
