package typictx

// AppModule is application module
type AppModule interface {
	GetMockTargets() []string
	GetTestTargets() []string
	GetConstructors() []interface{}
	GetCommands() []Command
	GetAction() Action
}
