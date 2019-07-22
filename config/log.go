package config

import (
	log "github.com/sirupsen/logrus"
)

var _ = func() bool {
	log.SetLevel(log.InfoLevel)
	// log.SetFormatter(&log.JSONFormatter{})

	return true
}()