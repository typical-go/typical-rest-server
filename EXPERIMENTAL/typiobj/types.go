package typiobj

import (
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Runner responsible to run the application
type Runner interface {
	Run(c *dig.Container) error
}

// Provider responsible to provide dependency
type Provider interface {
	Provide() []interface{}
}

// CommandLiner responsible to give command
type CommandLiner interface {
	CommandLine() cli.Command
}

// Configurer responsible to create config
type Configurer interface {
	Configure() Configuration
}

// Help model
type Help struct {
	// WIP:
	Name        string
	Description string
	// Configuration string
}
