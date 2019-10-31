package typiobj

import (
	"go.uber.org/dig"
)

// Provider responsible to provide dependency
type Provider interface {
	Provide() []interface{}
}

// Provide to execute provider
func Provide(c *dig.Container, p Provider) (err error) {
	for _, constructor := range p.Provide() {
		if err = c.Provide(constructor); err != nil {
			return
		}
	}
	return
}

// IsProvider return true if object implementation of provider
func IsProvider(obj interface{}) (ok bool) {
	_, ok = obj.(Provider)
	return
}
