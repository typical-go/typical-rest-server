package typiobj

// Destroyer responsible to destruct dependency
type Destroyer interface {
	Destroy() []interface{}
}

// IsDestroyer return true if object implementation of destructor
func IsDestroyer(obj interface{}) (ok bool) {
	_, ok = obj.(Destroyer)
	return
}
