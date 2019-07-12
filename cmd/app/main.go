package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/experimental/typiapp"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	app := typiapp.NewTypicalApplication(typical.Context)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
