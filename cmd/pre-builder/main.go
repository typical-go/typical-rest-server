package main

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	prebuilder.Prebuild(typical.Context)
}
