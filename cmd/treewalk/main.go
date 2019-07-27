package main

import (
	"fmt"
	"log"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
)

func main() {
	autowireFuncs, automockFiles, err := typiparser.Parse("app")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(autowireFuncs)
	fmt.Println(automockFiles)
}
