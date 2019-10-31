package typiobj

import "go.uber.org/dig"

// Destroyer responsible to destruct dependency
type Destroyer interface {
	Destroy() []interface{}
}

// Destroy to execute destroyer
func Destroy(c *dig.Container, d Destroyer) (err error) {
	for _, destructor := range d.Destroy() {
		if err = c.Invoke(destructor); err != nil {
			return
		}
	}
	return
}

// IsDestroyer return true if object implementation of destructor
func IsDestroyer(obj interface{}) (ok bool) {
	_, ok = obj.(Destroyer)
	return
}
