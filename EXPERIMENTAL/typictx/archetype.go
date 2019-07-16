package typictx

type ArcheType interface {
	GetMockTargets() []string
	GetTestTargets() []string
	GetConstructors() []interface{}
	GetConfig() interface{}
	GetConfigPrefix() string
	GetCommands() []Command
	GetAction() Action
}
