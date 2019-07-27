package main

import (
	"fmt"
	"log"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiparser"
)

func main() {
	paths := []string{"app", "app/controller", "app/helper", "app/repository", "app/service"}
	autowireFuncs, automockFiles, err := typiparser.Parse(paths)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(autowireFuncs)
	fmt.Println(automockFiles)
}
