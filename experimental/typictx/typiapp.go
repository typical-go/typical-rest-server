package typictx

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

// ConfigPrefixOrDefault return config prefix of typiapp or default value
func (a TypiApp) ConfigPrefixOrDefault() string {
	if a.ConfigPrefix == "" {
		return defaultConfigPrefix
	}
	return a.ConfigPrefix
}
