// +build typical

package main

import (
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	_ "github.com/typical-go/typical-rest-server/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typbuildtool.Run(typical.Context)
}
