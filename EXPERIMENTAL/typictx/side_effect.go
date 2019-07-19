package typictx

// SideEffect represent side effect from import
type SideEffect struct {
	Library            string
	AppFlag            bool
	TypicalDevToolFlag bool
}

// NewSideEffect return new instance of SideEffect
func NewSideEffect(library string) *SideEffect {
	return &SideEffect{
		Library:            library,
		AppFlag:            true,
		TypicalDevToolFlag: true,
	}
}

// ExcludeApp to exclude side effect from application import
func (e *SideEffect) ExcludeApp() *SideEffect {
	e.AppFlag = false
	return e
}

// ExcludeTypicalDevTool to exclude side effect from typical dev tool import
func (e *SideEffect) ExcludeTypicalDevTool() *SideEffect {
	e.TypicalDevToolFlag = false
	return e
}
