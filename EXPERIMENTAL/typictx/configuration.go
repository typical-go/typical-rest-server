package typictx

// Configuration represent the configuration
type Configuration struct {
	Prefix string
	Spec   interface{}
}

// Configure return config itself
func (c Configuration) Configure() Configuration {
	return c
}
