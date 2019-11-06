package prebuilder

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/collection"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	ProjectFiles       *walker.ProjectFiles
	Dirs               collection.Strings
	ApplicationImports golang.Imports
	ContextImport      string
	ConfigFields       []typiobj.ConfigField
	BuildCommands      []string
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	var files collection.Strings
	if p.Dirs, files, err = scanProject(typienv.AppName); err != nil {
		return
	}
	if p.ProjectFiles, err = walker.WalkProject(files); err != nil {
		return
	}
	p.ContextImport = ctx.Root + "/typical"
	log.Debug("Create imports for Application")
	for _, dir := range p.Dirs {
		p.ApplicationImports.AddImport("", ctx.Root+"/"+dir)
	}
	p.ApplicationImports.AddImport("", p.ContextImport)
	p.ConfigFields = typictx.ConfigFields(ctx)
	for _, command := range buildtool.Commands(ctx) {
		for _, subcommand := range command.Subcommands {
			s := fmt.Sprintf("%s_%s", command.Name, subcommand.Name)
			p.BuildCommands = append(p.BuildCommands, s)
		}
	}
	return
}
