package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	preBuilder := typicmd.NewTypicalPreBuilder(typical.Context)
	err := preBuilder.Cli().Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
