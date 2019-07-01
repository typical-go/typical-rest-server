package typienv

import "os"

const (
	envBin      = "TYPICAL_BIN"
	envCmd      = "TYPICAL_CMD"
	envName     = "TYPICAL_NAME"
	defaultBin  = "bin"
	defaultCmd  = "cmd"
	defaultName = "typical"
)

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
