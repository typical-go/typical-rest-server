package typictx

// TypiApp contain information of Typical Application
type TypiApp struct {
	Constructors []interface{}
	Action       Action
	Commands     []Command

	ConfigPrefix   string
	Config         interface{}
	ConfigLoadFunc interface{}

	MockTargets []string
	TestTargets []string
}

// GetConfig to get config
func (a TypiApp) GetConfig() interface{} {
	return a.Config
}

// GetConfigPrefix to get config
func (a TypiApp) GetConfigPrefix() string {
	return a.ConfigPrefix
}

// GetMockTargets to get mock targets
func (a TypiApp) GetMockTargets() []string {
	return a.MockTargets
}

// GetTestTargets to get test targets
func (a TypiApp) GetTestTargets() []string {
	return a.TestTargets
}

// GetConstructors to get constructors
func (a TypiApp) GetConstructors() []interface{} {
	constructors := a.Constructors
	constructors = append(constructors, a.ConfigLoadFunc)
	return constructors
}

// GetCommands to get commands
func (a TypiApp) GetCommands() []Command {
	return a.Commands
}

// GetAction to get commands
func (a TypiApp) GetAction() Action {
	return a.Action
}
