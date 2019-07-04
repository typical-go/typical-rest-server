package appctx

// TypiApp is Typical Application
type TypiApp struct {
	ConfigLoader
	Constructors   []interface{}
	Action         interface{}
	BinaryName     string
	ApplicationPkg string
	MockPkg        string
	MockTargets    []string
	TestTargets    []string
}
