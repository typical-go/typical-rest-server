package appctx

import "gopkg.in/urfave/cli.v1"

// Context of typical application
type Context struct {
	Name           string
	ConfigLoader   ConfigLoader
	Version        string
	Description    string
	Constructors   []interface{}
	Commands       []cli.Command
	Modules        []Module
	ReadmeTemplate string
}
