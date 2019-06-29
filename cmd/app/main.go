package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/typical"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = typical.Context.Name
	app.Usage = ""
	app.Description = typical.Context.Description
	app.Version = typical.Context.Version
	app.Action = actionFunction

	app.Run(os.Args)
}

func actionFunction(ctx *cli.Context) error {
	return typical.Context.Container().Invoke(func(s *app.Server) error {
		return s.Serve()
	})
}
