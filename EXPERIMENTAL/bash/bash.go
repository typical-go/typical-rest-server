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
