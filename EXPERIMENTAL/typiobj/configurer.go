package typiobj

// Configurer responsible to create config
type Configurer interface {
	Configure() Configuration
}

// IsConfigurer return true if object implementation of configurer
func IsConfigurer(obj interface{}) (ok bool) {
	_, ok = obj.(Configurer)
	return
}

// Configuration represent the configuration
type Configuration struct {
	Prefix string
	Spec   interface{}
}

// Configure return config itself
func (c Configuration) Configure() Configuration {
	return c
}
