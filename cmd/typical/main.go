package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical"
	"github.com/typical-go/typical-rest-server/typical/extension"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
)

func main() {
	typi := NewTypical(typical.Context)
	typi.AddExtension(&extension.GoExtension{})
	typi.AddExtension(&extension.ProjectExtension{})
	typi.AddExtension(&extension.DatabaseExtension{})

	err := typi.RunCLI(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
