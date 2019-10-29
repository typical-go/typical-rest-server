package typictx

import (
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Constructor responsible to construct dependency
type Constructor interface {
	Construct(c *dig.Container) error
}

// Destructor responsible to destruct dependency
type Destructor interface {
	Destruct(c *dig.Container) error
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
