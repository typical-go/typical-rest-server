package typiobj

import (
	"go.uber.org/dig"
)

// Preparer responsible to prepare
type Preparer interface {
	Prepare() []interface{}
}

// IsPreparer return true obj implement Preparer
func IsPreparer(obj interface{}) (ok bool) {
	_, ok = obj.(Preparer)
	return
}

// Prepare to execute preparer
func Prepare(c *dig.Container, preparer Preparer) (err error) {
	for _, preparation := range preparer.Prepare() {
		if err = c.Invoke(preparation); err != nil {
			return
		}
	}
	return
}
