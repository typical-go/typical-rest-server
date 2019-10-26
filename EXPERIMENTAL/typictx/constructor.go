package typictx

import "go.uber.org/dig"

// Constructor resposible to construct/destruct dependency
type Constructor interface {
	Construct(c *dig.Container) error
	Destruct(c *dig.Container) error
}
