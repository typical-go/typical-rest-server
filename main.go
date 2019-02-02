package main

import (
	"log"
	"os"

	"github.com/imantung/typical-go-server/app"
)

func main() {
	app := app.New()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
