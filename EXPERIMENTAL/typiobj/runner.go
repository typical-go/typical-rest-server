package typiobj

import (
	"go.uber.org/dig"
)

// Runner responsible to run the application
type Runner interface {
	Run(c *dig.Container) error
}
