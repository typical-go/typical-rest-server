package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical"
	"github.com/typical-go/typical-rest-server/typical/ext/xdb"
	"github.com/typical-go/typical-rest-server/typical/ext/xgo"
	"github.com/typical-go/typical-rest-server/typical/ext/xproj"
	"github.com/typical-go/typical-rest-server/typical/typigo"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
)

func main() {
	typi := typigo.NewTypical(typical.Context)
	typi.AddExtension(&xgo.GoExtension{})
	typi.AddExtension(xproj.NewProjectExtension(typical.Context))
	typi.AddExtension(&xdb.DatabaseExtension{})

	err := typi.RunCLI(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
