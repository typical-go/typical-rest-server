package main

import (
	"github.com/typical-go/typical-go/pkg/typicmd/application"
	_ "github.com/typical-go/typical-rest-server/cmd/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	application.Run(typical.Context)
}
