package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/typical-go/typical-rest-server/typical"
	"github.com/typical-go/typical-rest-server/typical/typiapp"
)

func main() {
	app := typiapp.NewTypicalApplication(typical.Context)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
