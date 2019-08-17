package typserver

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module() *typictx.Module {
	return &typictx.Module{
		Name:      "Echo Server with Logrus",
		OpenFunc:  Create,
		CloseFunc: Shutdown,

		Config: typictx.Config{
			Prefix: "SERVER",
			Spec:   &Config{},
		},
	}
}
