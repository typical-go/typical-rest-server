package appctx

import (
	"go.uber.org/dig"
)

// Context of typical application
type Context struct {
	Name           string
	ConfigPrefix   string
	Path           string
	Version        string
	Description    string
	Container      *dig.Container
	Modules        []interface{}
	ReadmeTemplate string
}

// ConfigDoc for configuration documentation
func (c Context) ConfigDoc() string {
	return "mehmehmeh"
}
