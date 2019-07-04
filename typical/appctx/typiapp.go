package appctx

// TypiApp contain information of Typical Application
type TypiApp struct {
	ConfigLoader
	Constructors   []interface{}
	Action         interface{}
	Commands       []Command
	BinaryName     string
	ApplicationPkg string
	MockPkg        string
	MockTargets    []string
	TestTargets    []string
}
