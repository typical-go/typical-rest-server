package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool"
	_ "github.com/typical-go/typical-rest-server/cmd/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	err := buildtool.Cli(typical.Context).Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
