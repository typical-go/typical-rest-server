package debugkit

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// ElapsedTime to print elapsed time of function
func ElapsedTime(name string) func() {
	start := time.Now()
	return func() {
		log.Debugf("%s took %v\n", name, time.Since(start))
	}
}
