package appctx

// Context of typical application
type Context struct {
	TypiApp
	TypiCli

	Name        string
	Version     string
	Description string

	Modules        []Module
	ReadmeTemplate string
}
