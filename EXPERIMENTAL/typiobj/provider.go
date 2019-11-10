package typiobj

// Provider responsible to provide dependency
type Provider interface {
	Provide() []interface{}
}

// IsProvider return true if object implementation of provider
func IsProvider(obj interface{}) (ok bool) {
	_, ok = obj.(Provider)
	return
}
