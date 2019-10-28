package typictx

// Configurer responsible to create config
type Configurer interface {
	Configure() Config
}

// Config represent the configuration
type Config struct {
	Prefix string
	Spec   interface{}
}
