package typiobj

// Runner responsible to run the application
type Runner interface {
	Run() interface{}
}

// IsRunner return true if obj implement Runner
func IsRunner(obj interface{}) (ok bool) {
	_, ok = obj.(Runner)
	return
}
