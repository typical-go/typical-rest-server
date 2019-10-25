package prebuilder

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	MockTarget    *MockTargetGenerator
	Constructor   *ConstructorGenerator
	Configuration *ConfigurationGenerator
	TestTarget    *TestTargetGenerator
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	var projectFiles *walker.ProjectFiles
	var contextFile *walker.ContextFile
	var dirs, files []string
	log.Debug("Scan project to get package and filenames")
	if dirs, files, err = scanProject(typienv.AppName); err != nil {
		return
	}
	log.Debug("Walk the project to get annotated or metadata")
	if projectFiles, err = walker.WalkProject(files); err != nil {
		return
	}
	log.Debug("Walk the context file")
	if contextFile, err = walker.WalkContext(ctxPath); err != nil {
		return
	}
	p.MockTarget = &MockTargetGenerator{
		Root:        ctx.Root,
		MockTargets: projectFiles.Automocks(),
		Packages:    dirs,
	}
	p.Constructor = &ConstructorGenerator{
		Root:         ctx.Root,
		Constructors: projectFiles.Autowires(),
		Packages:     dirs,
	}
	p.Configuration = &ConfigurationGenerator{
		Configs:     createConfigs(ctx),
		ContextFile: contextFile,
	}
	p.TestTarget = &TestTargetGenerator{
		Root:     ctx.Root,
		Packages: packages,
	}
	return
}

func (p *prebuilder) Prebuild() (r report, err error) {
	if r.TestTargetUpdated, err = p.TestTarget.Generate(); err != nil {
		return
	}
	if r.MockTargetUpdated, err = p.MockTarget.Generate(); err != nil {
		return
	}
	if r.ConstructorUpdated, err = p.Constructor.Generate(); err != nil {
		return
	}
	if r.ConfigurationUpdated, err = p.Configuration.Generate(); err != nil {
		return
	}
	return
}

func createConfigs(ctx *typictx.Context) (configs []Config) {
	configs = append(configs, Config{
		Key: strcase.ToCamel(strings.ToLower(ctx.Application.Prefix)),
		Typ: reflect.TypeOf(ctx.Application.Spec).String(),
	})
	for _, module := range ctx.Modules {
		configs = append(configs, Config{
			Key: strcase.ToCamel(strings.ToLower(module.Prefix)),
			Typ: reflect.TypeOf(module.Spec).String(),
		})
	}
	return
}
