package prebuilder

import (
	log "github.com/sirupsen/logrus"

	"time"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

type prebuilder struct {
	// *typictx.Context
	// *walker.ProjectFiles
	// *walker.ContextFile
	// Filenames []string
	// Packages  []string
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
		Context:      ctx,
		ProjectFiles: projectFiles,
		Packages:     packages,
	}
	p.Configuration = &ConfigurationGenerator{
		Context:     ctx,
		ContextFile: contextFile,
	}
	p.TestTarget = &TestTargetGenerator{
		Context:  ctx,
		Packages: packages,
	}

	return
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Debugf("%s took %v\n", what, time.Since(start))
	}
}
