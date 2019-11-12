package typienv

import (
	"os"
)

// TODO: create struct for envkey and default value
const (
	envApp        = "TYPICAL_APP"
	envBin        = "TYPICAL_BIN"
	envCmd        = "TYPICAL_CMD"
	envMock       = "TYPICAL_MOCK"
	envRelease    = "TYPICAL_RELEASE"
	envBuildTool  = "TYPICAL_BUILD_TOOL"
	envDependency = "TYPICAL_DEPENDENCY"

	defaultBin        = "bin"
	defaultCmd        = "cmd"
	defaultMock       = "mock"
	defaultApp        = "app"
	defaultRelease    = "release"
	defaultBuildTool  = "build-tool"
	defaultDependency = "dependency"
	defaultMetadata   = ".typical-metadata"
	defaultReadme     = "README.md"
)

var (
	App        *applicationFolder
	BuildTool  *applicationFolder
	Dependency *applicationFolder
	Bin        string
	Metadata   string
	Mock       string
	Release    string
	AppName    string
	Readme     string
)

type applicationFolder struct {
	Package string
	SrcPath string
	BinPath string
}

func init() {
	AppName = app()
	cmd := cmd()
	Bin = bin()
	Metadata = defaultMetadata
	buildTool := buildTool()
	dependency := dependency()
	App = &applicationFolder{
		Package: "main",
		SrcPath: cmd + "/" + AppName,
		BinPath: Bin + "/" + AppName,
	}
	BuildTool = &applicationFolder{
		Package: "main",
		SrcPath: cmd + "/" + buildTool,
		BinPath: Bin + "/" + buildTool,
	}
	Dependency = &applicationFolder{
		Package: dependency,
		SrcPath: cmd + "/internal/" + dependency,
	}
	Mock = mock()
	Release = release()
	Readme = defaultReadme
}

func cmd() string {
	cmd := os.Getenv(envCmd)
	if cmd == "" {
		cmd = defaultCmd
	}
	return cmd
}

func bin() string {
	bin := os.Getenv(envBin)
	if bin == "" {
		bin = defaultBin
	}
	return bin
}

func buildTool() string {
	devTool := os.Getenv(envBuildTool)
	if devTool == "" {
		devTool = defaultBuildTool
	}
	return devTool
}

func app() string {
	app := os.Getenv(envApp)
	if app == "" {
		app = defaultApp
	}
	return app
}

func mock() string {
	mock := os.Getenv(envMock)
	if mock == "" {
		mock = defaultMock
	}
	return mock
}

func release() string {
	release := os.Getenv(envApp)
	if release == "" {
		release = defaultRelease
	}
	return release
}

func dependency() string {
	dep := os.Getenv(envDependency)
	if dep == "" {
		dep = defaultDependency
	}
	return dep
}
