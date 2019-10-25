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
	log.Debug("BinaryNotExist: ", c.BinaryNotExist, " ",
		"PrebuildUpdated: ", c.PrebuildUpdated, " ",
		"HaveBuildArgs: ", c.HaveBuildArgs)
	return c.BinaryNotExist ||
		c.PrebuildUpdated ||
		c.HaveBuildArgs
}
