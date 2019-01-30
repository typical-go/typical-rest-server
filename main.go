package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	AppName  = "Typical Go Server"
	AppUsage = "API Server for [purpose]"
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
