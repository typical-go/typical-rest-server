package typictx

// Config represent the configuration
type Config interface {
	Prefix() string
	Spec() interface{}
}

// NewConfig return new instance of Config
func NewConfig(prefix string, spec interface{}) Config {
	return config{
		prefix: prefix,
		spec:   spec,
	}
}

type config struct {
	prefix string
	spec   interface{}
}

func (c config) Prefix() string {
	return c.prefix
}

func (c config) Spec() interface{} {
	return c.spec
}
