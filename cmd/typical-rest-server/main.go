package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-rest-server/internal/app"
	_ "github.com/typical-go/typical-rest-server/internal/generated"
)

func main() {
	typapp.Start(app.Start)
}
