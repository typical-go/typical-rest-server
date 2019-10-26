package prebuilder

import (
	log "github.com/sirupsen/logrus"
)

type buildToolChecker struct {
	BinaryNotExist  bool
	PrebuildUpdated bool
	HaveBuildArgs   bool
}

func (c *buildToolChecker) Check() bool {
	log.WithFields(log.Fields{
		"BinaryNotExist":  c.BinaryNotExist,
		"PrebuildUpdated": c.PrebuildUpdated,
		"HaveBuildArgs":   c.HaveBuildArgs,
	}).Debug("Check for build-tool")
	return c.BinaryNotExist ||
		c.PrebuildUpdated ||
		c.HaveBuildArgs
}
