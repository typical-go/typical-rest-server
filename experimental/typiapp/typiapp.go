package typiapp

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/typical-go/typical-rest-server/experimental/typictx"
	"gopkg.in/urfave/cli.v1"
)

// TypicalApplication represent typical application
type TypicalApplication struct {
	typictx.Context
}

// NewTypicalApplication return new instance of TypicalApplications
func NewTypicalApplication(context typictx.Context) *TypicalApplication {
	return &TypicalApplication{context}
}

// Run the typical application
func (t *TypicalApplication) Run(arguments []string) error {
	app := cli.NewApp()
	app.Name = t.Name
	app.Usage = ""
	app.Description = t.Description
	app.Version = t.Version
	app.Action = t.startApplication

	for i := range t.TypiApp.Commands {
		cmd := t.TypiApp.Commands[i]
		app.Commands = append(app.Commands, cli.Command{
			Name:      cmd.Name,
			ShortName: cmd.ShortName,
			Action:    t.invoke(cmd.Action),
		})
	}

	return app.Run(arguments)
}

func (t *TypicalApplication) startApplication(ctx *cli.Context) {
	container := t.Container()

	if t.StopFunc != nil {
		gracefulStop := make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		// gracefull shutdown
		go func() {
			<-gracefulStop
			container.Invoke(t.StopFunc)
		}()
	}

	if t.StartFunc != nil {
		container.Invoke(t.StartFunc)
	}
}

func (t *TypicalApplication) invoke(f interface{}) interface{} {
	return func(ctx *cli.Context) error {
		return t.Container().Invoke(f)
	}
}
