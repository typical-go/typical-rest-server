package appctx

// Context of typical application
type Context struct {
	Name           string
	ConfigPrefix   string
	Path           string
	Version        string
	Description    string
	Config         interface{}
	Constructors   []interface{}
	Modules        map[string]Module
	ReadmeTemplate string
}
