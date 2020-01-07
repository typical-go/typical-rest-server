package main

import (
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typbuildtool.Run(&typical.Descriptor)
}
