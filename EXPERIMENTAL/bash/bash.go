package bash

import (
	"os"
	"os/exec"
)

// Run to execute bash or exit the application
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
