package appctx

// TypiApp contain information of Typical Application
type TypiApp struct {
	Constructors []interface{}
	Action       interface{}
	Commands     []Command

	// configuration related
	ConfigPrefix   string
	Config         interface{}
	ConfigLoadFunc interface{}

	// build related
	BinaryName     string
	ApplicationPkg string
	MockPkg        string
	MockTargets    []string
	TestTargets    []string
}
