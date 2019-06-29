package appctx

// Context of typical application
type Context struct {
	Name           string
	ConfigPrefix   string
	Path           string
	Version        string
	Description    string
	Constructors   []interface{}
	Modules        map[string]Module
	ReadmeTemplate string
}
