package main

import (
	
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/typical"
	"github.com/typical-go/typical-rest-server/experimental/typicli"
	
)

func main() {
	cli := typicli.NewTypicalCli(typical.Context)
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
