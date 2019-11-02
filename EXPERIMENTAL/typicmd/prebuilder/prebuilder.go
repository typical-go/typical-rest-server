package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	ProjectFiles       *walker.ProjectFiles
	Dirs               slice.Strings
	ApplicationImports golang.Imports
	ContextImport      string
	ConfigFields       []typiobj.ConfigField
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	var files slice.Strings
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
	return
}
