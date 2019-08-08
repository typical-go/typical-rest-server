package bash

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Run to execute bash or exit the application
func Run(name string, args ...string) error {
	log.Infof("Run: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return cmd.Run()
}

// Silent same with `Run()` without print any output
func Silent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
