package typictx

// TypiApp contain information of Typical Application
type TypiApp struct {
	Constructors []interface{}
	Action       Action
	Commands     []Command
	MockTargets  []string
	TestTargets  []string
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
	return a.Constructors
}

// GetCommands to get commands
func (a TypiApp) GetCommands() []Command {
	return a.Commands
}

// GetAction to get commands
func (a TypiApp) GetAction() Action {
	return a.Action
}
