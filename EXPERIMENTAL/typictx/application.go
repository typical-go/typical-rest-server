package typictx

// Application is represent the application
type Application struct {
	Config      Config
	StartFunc   interface{}
	StopFunc    interface{}
	Commands    []*Command
	Initiations []interface{}
}

// Configure return configuration
func (a Application) Configure() Config {
	return a.Config
}
