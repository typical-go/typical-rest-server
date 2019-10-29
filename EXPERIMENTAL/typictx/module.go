package typictx

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"
	"github.com/typical-go/typical-rest-server/pkg/utility/errkit"
	"go.uber.org/dig"
)

// Module of typical-go application
type Module struct {
	Configuration
	Name         string
	OpenFunc     interface{} // TODO: remove this, use constructors
	CloseFunc    interface{} // TODO: remove this, use destructors
	Constructors slice.Interfaces
	Destructors  slice.Interfaces
}

// Construct dependency
func (m Module) Construct(c *dig.Container) (err error) {
	for _, constructor := range m.Constructors {
		if err = c.Provide(constructor); err != nil {
			return
		}
	}
	return c.Provide(m.OpenFunc)
}

// Destruct dependency
func (m Module) Destruct(c *dig.Container) (err error) {
	var errs errkit.Errors
	for _, destructor := range m.Destructors {
		errs.Add(c.Invoke(destructor))
	}
	errs.Add(c.Invoke(m.CloseFunc))
	return errs
}
