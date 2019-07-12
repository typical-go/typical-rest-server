package typictx

import (
	"os"
	"os/signal"
	"syscall"
)

// TypiApp contain information of Typical Application
type TypiApp struct {
	StartFunc interface{}
	StopFunc  interface{}

	Constructors []interface{}
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

// StartApplication to start the application
func (a TypiApp) StartApplication(ctx StartContext) {
	container := ctx.Container()

	if a.StopFunc != nil {
		gracefulStop := make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		// gracefull shutdown
		go func() {
			<-gracefulStop
			container.Invoke(a.StopFunc)
		}()
	}

	if a.StartFunc != nil {
		container.Invoke(a.StartFunc)
	}
}
