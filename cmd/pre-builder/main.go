// +build typical

package main

import (
	"github.com/typical-go/typical-go/pkg/typprebuilder"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typprebuilder.Run(typical.Context)
}
