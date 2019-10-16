package prebuilder

import (
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	Annotated     *AnnotatedGenerator
	Configuration *ConfigurationGenerator
	TestTarget    *TestTargetGenerator
}

func (p *prebuilder) Initiate(ctx *typictx.Context) (err error) {
	log.Debug("Scan project to get package and filenames")
	root := typienv.AppName
	packages, filenames, err := scanProject(root)
	if err != nil {
		return
	}
	log.Debug("Walk the project to get annotated or metadata")
	projectFiles, err := walker.WalkProject(filenames)
	if err != nil {
		return
	}
	log.Debug("Walk the context file")
	contextFile, err := walker.WalkContext(ctxPath)
	if err != nil {
		return
	}
	p.Annotated = &AnnotatedGenerator{
		Root:         ctx.Root,
		ProjectFiles: projectFiles,
		Packages:     packages,
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

func (p *prebuilder) Prebuild() (err error) {
	if err = p.TestTarget.Generate(); err != nil {
		return
	}
	if err = p.Annotated.Generate(); err != nil {
		return
	}
	if err = p.Configuration.Generate(); err != nil {
		return
	}
	return
}

func createConfigs(ctx *typictx.Context) (configs []Config) {
	for _, acc := range ctx.ConfigAccessors() {
		key := acc.GetKey()
		typ := reflect.TypeOf(acc.GetConfigSpec()).String()
		configs = append(configs, Config{
			Key: key,
			Typ: typ,
		})
	}
	return
}
