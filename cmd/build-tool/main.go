package main

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool"
	_ "github.com/typical-go/typical-rest-server/cmd/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	buildtool.Run(typical.Context)
}
