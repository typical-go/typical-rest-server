package bash

import (
	"log"
	"os"
	"os/exec"
)

// Run to execute bash or exit the application
func Run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Silent same with `Run()` without print any output
func Silent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
