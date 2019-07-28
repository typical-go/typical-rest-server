package main

import (
	"fmt"
	"log"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
)

func main() {
	projCtx, err := typiparser.Parse("app")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v", projCtx)

}
