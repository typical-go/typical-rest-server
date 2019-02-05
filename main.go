package main

import (
	"log"

	"github.com/imantung/typical-go-server/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
