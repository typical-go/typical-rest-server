package typienv

import (
	"os"
)

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
)

var (
	App        *applicationFolder
	BuildTool  *applicationFolder
	Mock       string
	Release    string
	Dependency string
)

type applicationFolder struct {
	MainPkg string
	Binary  string
}

func init() {
	cmd := cmd()
	bin := bin()
	app := app()
	buildTool := buildTool()
	App = &applicationFolder{
		MainPkg: cmd + "/" + app,
		Binary:  bin + "/" + app,
	}
	BuildTool = &applicationFolder{
		MainPkg: cmd + "/" + buildTool,
		Binary:  bin + "/" + buildTool,
	}
	Mock = mock()
	Release = release()
	Dependency = cmd + "/internal/" + dependency()
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
