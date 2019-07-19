package typictx

// AppModule is application module
type AppModule interface {
	GetMockTargets() []string
	GetTestTargets() []string
	GetConstructors() []interface{}
	GetConfig() interface{}
	GetConfigPrefix() string
	GetCommands() []Command
	GetAction() Action
}
