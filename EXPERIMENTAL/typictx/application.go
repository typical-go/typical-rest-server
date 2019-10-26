package typictx

import "github.com/urfave/cli"

// Application is represent the application
type Application struct {
	Config      Config
	StartFunc   interface{}
	StopFunc    interface{}
	Commands    []cli.Command
	Initiations []interface{}
}

// Configure return configuration
func (a Application) Configure() Config {
	return a.Config
}
