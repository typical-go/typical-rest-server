package main

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typcore.RunApp(&typical.Descriptor)
}
