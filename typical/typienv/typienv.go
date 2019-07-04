package typienv

import (
	"fmt"
	"os"
)

const (
	envBin      = "TYPICAL_BIN"
	envCmd      = "TYPICAL_CMD"
	envName     = "TYPICAL_NAME"
	defaultBin  = "bin"
	defaultCmd  = "cmd"
	defaultName = "typical"
)

// BinaryPath return complete path of typical binary
func BinaryPath(name string) string {
	return fmt.Sprintf("./%s/%s", Bin(), name)
}

// MainPackage return main package path
func MainPackage(name string) string {
	return fmt.Sprintf("./%s/%s", Cmd(), name)
}

// TypicalMainPackage return main package path of Typical CLI
func TypicalMainPackage() string {
	return MainPackage(Name())
}

// TypicalBinaryPath return
func TypicalBinaryPath() string {
	return BinaryPath(Name())
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

// Name to return typical cmd name
func Name() string {
	name := os.Getenv(envName)
	if name == "" {
		name = defaultName
	}
	return name
}
