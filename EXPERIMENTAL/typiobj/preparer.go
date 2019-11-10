package typiobj

// Preparer responsible to prepare
type Preparer interface {
	Prepare() []interface{}
}

// IsPreparer return true obj implement Preparer
func IsPreparer(obj interface{}) (ok bool) {
	_, ok = obj.(Preparer)
	return
}
