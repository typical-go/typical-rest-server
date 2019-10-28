package prebuilder

import (
	"reflect"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/golang"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/slice"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	ProjectFiles       *walker.ProjectFiles
	Dirs               slice.Strings
	ConfigImports      golang.Imports
	ApplicationImports golang.Imports
	ContextImport      string
	Configs            []config
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	var contextFile *walker.ContextFile
	var files slice.Strings
	log.Debug("Scan project to get package and filenames")
	if p.Dirs, files, err = scanProject(typienv.AppName); err != nil {
		return
	}
	log.Debug("Walk the project to get annotated or metadata")
	if p.ProjectFiles, err = walker.WalkProject(files); err != nil {
		return
	}
	log.Debug("Walk the context file")
	if contextFile, err = walker.WalkContext(ctxPath); err != nil {
		return
	}
	log.Debug("Create context import")
	p.ContextImport = ctx.Root + "/typical"
	log.Debug("Create imports for Config")
	p.ConfigImports = contextFile.Imports
	p.ConfigImports.AddImport("", "github.com/kelseyhightower/envconfig")
	log.Debug("Create imports for Application")
	for _, dir := range p.Dirs {
		p.ApplicationImports.AddImport("", ctx.Root+"/"+dir)
	}
	p.ApplicationImports.AddImport("", p.ContextImport)
	log.Debug("Create configs")
	for _, cfg := range ctx.Configs() {
		p.Configs = append(p.Configs, config{
			Key: fmtConfigKey(cfg.Prefix),
			Typ: fmtConfigTyp(cfg.Spec),
		})
	}
	return
}

func (p *prebuilder) Prebuild() (r report, err error) {
	if r.TestTargetUpdated, err = Generate("test_target", testTarget{
		ContextImport: p.ContextImport,
		Packages:      p.Dirs,
	}); err != nil {
		return
	}
	if r.MockTargetUpdated, err = Generate("mock_target", mockTarget{
		ApplicationImports: p.ApplicationImports,
		MockTargets:        p.ProjectFiles.Automocks(),
	}); err != nil {
		return
	}
	if r.ConstructorUpdated, err = Generate("constructor", constructor{
		ApplicationImports: p.ApplicationImports,
		Constructors:       p.ProjectFiles.Autowires(),
	}); err != nil {
		return
	}
	if r.ConfigurationUpdated, err = Generate("configuration", configuration{
		Configs:       p.Configs,
		ConfigImports: p.ConfigImports,
	}); err != nil {
		return
	}
	return
}

func fmtConfigKey(s string) string {
	return strcase.ToCamel(strings.ToLower(s))
}

func fmtConfigTyp(v interface{}) string {
	return reflect.TypeOf(v).String()
}
