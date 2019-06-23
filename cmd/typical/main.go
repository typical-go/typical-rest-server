package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical"
	"github.com/typical-go/typical-rest-server/typical/ext/xdb"
	"github.com/typical-go/typical-rest-server/typical/ext/xgo"
	"github.com/typical-go/typical-rest-server/typical/ext/xproj"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
)

func main() {
	typi := NewTypical(typical.Context)
	typi.AddExtension(&xgo.GoExtension{})
	typi.AddExtension(&xproj.ProjectExtension{})
	typi.AddExtension(&xdb.DatabaseExtension{})

	err := typi.RunCLI(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
