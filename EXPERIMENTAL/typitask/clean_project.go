package typitask

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// CleanProject to remove file from build processs
func CleanProject(ctx typictx.ActionContext) error {
	log.Info("Remove bin folder")
	os.RemoveAll(typienv.Bin())

	log.Info("Go clean")
	os.Setenv("GO111MODULE", "off") // NOTE:XXX: https://github.com/golang/go/issues/28680
	return bash.GoClean("-x", "-testcache", "-modcachœœe")
}
