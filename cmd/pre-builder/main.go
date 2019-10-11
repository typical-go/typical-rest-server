package main

import (
	"log"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	err := prebuilder.Prebuild(typical.Context)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
