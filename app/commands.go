package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"

	"github.com/urfave/cli"
)

func initCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{Name: "Serve", ShortName: "s", Usage: "Serve the clients", Action: cmdServe},
		// add more command here
	}
}

func cmdServe(c *cli.Context) error {
	e := echo.New()
	initMiddlewares(e)
	initRoutes(e)

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// gracefull shutdown
	go func() {
		<-gracefulStop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		e.Shutdown(ctx)
	}()

	return e.Start(conf.Address)

}
