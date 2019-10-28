package typserver

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module() *typictx.Module {
	return &typictx.Module{
		Name:      "Server",
		OpenFunc:  Create,
		CloseFunc: Shutdown,
		Config:    typictx.Config{"SERVER", &Config{}},
	}
}
