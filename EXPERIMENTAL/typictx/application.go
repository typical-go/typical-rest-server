package typictx

import "github.com/urfave/cli"

// Application is represent the application
type Application struct {
	Configuration
	StartFunc   interface{}
	StopFunc    interface{}
	Commands    []cli.Command
	Initiations []interface{}
}
