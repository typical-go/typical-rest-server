package typserver

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module() *typictx.Module {
	return &typictx.Module{
		Name:         "Echo Server with Logrus",
		ConfigPrefix: "SERVER",
		ConfigSpec:   &Config{},
		OpenFunc:     Create,
		CloseFunc:    Shutdown,
	}
}
