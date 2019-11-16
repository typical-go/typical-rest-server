package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	_ "github.com/typical-go/typical-rest-server/internal/dependency"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typapp.Run(typical.Context)
}
