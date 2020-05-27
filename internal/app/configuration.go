package app

import (
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app/config"
)

var (
	configName = "APP"
)

// Configuration of server
func Configuration() *typgo.Configuration {
	return &typgo.Configuration{
		Name: configName,
		Spec: &config.Config{
			Debug: true,
		},
	}
}
