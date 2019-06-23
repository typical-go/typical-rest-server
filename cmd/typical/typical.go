package main

import (
	"github.com/typical-go/typical-go/appx"
	"github.com/typical-go/typical-rest-server/typical/extension"
	"gopkg.in/urfave/cli.v1"
)

// Typical represent the typical framework tool
type Typical struct {
	Context    appx.Context
	extensions []extension.Extension
}

// NewTypical return new instance of Typical CLI
func NewTypical(context appx.Context) *Typical {
	return &Typical{
		Context: context,
	}
}

// AddExtension will add extension to typical
func (t *Typical) AddExtension(extension extension.Extension) {
	t.extensions = append(t.extensions, extension)
}

// RunCLI start the command line interface
func (t *Typical) RunCLI(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.Context.Name
	app.Usage = ""
	app.Description = t.Context.Description
	app.Version = t.Context.Version
	app.Commands = t.commands()
	return app.Run(arguments)
}

func (t *Typical) commands() (cmds []cli.Command) {
	for _, extension := range t.extensions {
		cmds = append(cmds, extension.Command())
	}
	return cmds
}
