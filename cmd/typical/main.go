package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typimain"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	cli := typimain.NewTypicalCli(typical.Context)
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
