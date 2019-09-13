package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typimain"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	buildTool := typimain.NewTypicalBuildTool(typical.Context)
	err := buildTool.Cli().Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
