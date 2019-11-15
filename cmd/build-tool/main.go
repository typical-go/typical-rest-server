package main

import (
	"github.com/typical-go/typical-go/pkg/typicmd/buildtool"
	_ "github.com/typical-go/typical-rest-server/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	buildtool.Run(typical.Context)
}
