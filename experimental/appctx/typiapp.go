package appctx

import (
	"os"
	"strings"
)

// TypiApp contain information of Typical Application
type TypiApp struct {
	Constructors []interface{}
	Action       interface{}
	Commands     []Command

	ConfigPrefix   string
	Config         interface{}
	ConfigLoadFunc interface{}

	BinaryName     string
	ApplicationPkg string
	MockPkg        string

	MockTargets []string
	TestTargets []string
}

// BinaryNameOrDefault return binary name of typiapp or default value
func (a TypiApp) BinaryNameOrDefault() string {
	if a.BinaryName == "" {
		dir, _ := os.Getwd()
		chunks := strings.Split(dir, "/")
		return chunks[len(chunks)-1]
	}
	return a.BinaryName
}

// ConfigPrefixOrDefault return config prefix of typiapp or default value
func (a TypiApp) ConfigPrefixOrDefault() string {
	if a.ConfigPrefix == "" {
		return defaultConfigPrefix
	}
	return a.ConfigPrefix
}

// ApplicationPkgOrDefault return application package of typiapp or default value
func (a TypiApp) ApplicationPkgOrDefault() string {
	if a.ApplicationPkg == "" {
		return defaultApplicationPkg
	}
	return a.ApplicationPkg
}

// MockPkgOrDefault return mock package of typiapp or default value
func (a TypiApp) MockPkgOrDefault() string {
	if a.MockPkg == "" {
		return defaultMockPkg
	}
	return a.MockPkg
}
