package typienv

import (
	"fmt"
	"os"
)

const (
	envApp       = "TYPICAL_APP"
	envBin       = "TYPICAL_BIN"
	envCmd       = "TYPICAL_CMD"
	envMock      = "TYPICAL_MOCK"
	envRelease   = "TYPICAL_RELEASE"
	envBuildTool = "TYPICAL_BUILD_TOOL"

	defaultBin       = "bin"
	defaultCmd       = "cmd"
	defaultMock      = "mock"
	defaultApp       = "app"
	defaultRelease   = "release"
	defaultBuildTool = "build-tool"
)

// Binary return complete path of typical binary
func Binary(name string) string {
	return fmt.Sprintf("./%s/%s", Bin(), name)
}

// MainPackage return main package path
func MainPackage(name string) string {
	return fmt.Sprintf("./%s/%s", Cmd(), name)
}

// Bin to return typical bin folder
func Bin() string {
	bin := os.Getenv(envBin)
	if bin == "" {
		bin = defaultBin
	}
	return bin
}

// Cmd to return typical cmd folder
func Cmd() string {
	cmd := os.Getenv(envCmd)
	if cmd == "" {
		cmd = defaultCmd
	}
	return cmd
}

// BuildTool to return typical dev tool package
func BuildTool() string {
	devTool := os.Getenv(envBuildTool)
	if devTool == "" {
		devTool = defaultBuildTool
	}
	return devTool
}

// BuildToolMainPackage return main package path of Typical CLI
func BuildToolMainPackage() string {
	return MainPackage(BuildTool())
}

// BuildToolBinary return
func BuildToolBinary() string {
	return Binary(BuildTool())
}

// App to return app package
func App() string {
	app := os.Getenv(envApp)
	if app == "" {
		app = defaultApp
	}
	return app
}

// AppMainPackage return main package path of Typical CLI
func AppMainPackage() string {
	return MainPackage(App())
}

// AppBinary return
func AppBinary() string {
	return Binary(App())
}

// Mock to return app package
func Mock() string {
	mock := os.Getenv(envMock)
	if mock == "" {
		mock = defaultMock
	}
	return mock
}

// Release to return release package
func Release() string {
	release := os.Getenv(envApp)
	if release == "" {
		release = defaultRelease
	}
	return release
}
