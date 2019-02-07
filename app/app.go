package app

import (
	"os"

	"github.com/imantung/typical-go-server/config"
	"github.com/urfave/cli"
)

// FIXME: Application Name, Description, Config Prefix
const (
	appName         = "[Service Name]"
	appUsage        = "API Server for [Service Description]"
	appVersion      = "0.0.1"
	appConfigPrefix = "APP"
)

var (
	app  cli.App
	conf config.Config
)

// Run the service
func Run() (err error) {

	// prepare configuration
	conf, err = config.Load(appConfigPrefix)
	if err != nil {
		return
	}

	// command line interface
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Version = appVersion

	initCommands(app)

	return app.Run(os.Args)
}
