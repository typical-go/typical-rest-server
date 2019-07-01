package appctx

import "gopkg.in/urfave/cli.v1"

// Context of typical application
type Context struct {
	Name           string
	ConfigPrefix   string
	Path           string
	Version        string
	Description    string
	Config         interface{}
	Constructors   []interface{}
	Commands       []cli.Command
	Modules        []Module
	ReadmeTemplate string
}
