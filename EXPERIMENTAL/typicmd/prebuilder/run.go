package prebuilder

import (
	"os"
	"os/exec"

	"github.com/typical-go/typical-rest-server/pkg/utility/filekit"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

var (
	app        = typienv.App.SrcPath
	buildTool  = typienv.BuildTool.SrcPath
	dependency = typienv.Dependency.SrcPath
)

const (
	ctxPath  = "typical/context.go"
	debugEnv = "PREBUILDER_DEBUG"
)

// Run the prebuilder
func Run(ctx *typictx.Context) {
	var err error
	var preb prebuilder
	checker := checker{
		contextChecksum: contextChecksum(),
		buildToolBinary: !filekit.IsExist(typienv.BuildTool.BinPath),
		readmeFile:      !filekit.IsExist(typienv.Readme),
	}
	if os.Getenv(debugEnv) != "" {
		log.SetLevel(log.DebugLevel)
	}
	if err = ctx.Validate(); err != nil {
		log.Fatal(err.Error())
	}
	if err = typictx.GenerateEnvfile(ctx); err != nil {
		log.Fatal(err.Error())
	}
	if err := preb.Initiate(ctx); err != nil {
		log.Fatal(err.Error())
	}
	if checker.configuration, err = metadata.Update("config_fields", preb.ConfigFields); err != nil {
		log.Fatal(err.Error())
	}
	if checker.buildCommands, err = metadata.Update("build_commands", preb.BuildCommands); err != nil {
		log.Fatal(err.Error())
	}
	if checker.testTarget, err = Generate("test_target", testTarget{
		ContextImport: preb.ContextImport,
		Packages:      preb.Dirs,
	}); err != nil {
		log.Fatal(err.Error())
	}
	if checker.mockTarget, err = Generate("mock_target", mockTarget{
		ApplicationImports: preb.ApplicationImports,
		MockTargets:        preb.ProjectFiles.Automocks(),
	}); err != nil {
		log.Fatal(err.Error())
	}
	if checker.constructor, err = Generate("constructor", constructor{
		ApplicationImports: preb.ApplicationImports,
		Constructors:       preb.ProjectFiles.Autowires(),
	}); err != nil {
		log.Fatal(err.Error())
	}
	if checker.checkBuildTool() {
		log.Info("Build the build-tool")
		if err := bash.GoBuild(typienv.BuildTool.BinPath, typienv.BuildTool.SrcPath); err != nil {
			log.Fatal(err.Error())
		}
	}
	if checker.checkReadme() {
		log.Info("Generate readme")
		cmd := exec.Command(typienv.BuildTool.BinPath, "readme")
		if err := cmd.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func contextChecksum() bool {
	// NOTE: context checksum is passed by typicalw
	if len(os.Args) > 1 {
		return os.Args[1] == "1"
	}
	return false
}
