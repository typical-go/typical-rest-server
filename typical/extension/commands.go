package extension

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical"
	"gopkg.in/urfave/cli.v1"
)

func invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := typical.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}

func notImplement(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func print(f func() string) interface{} {
	return func(ctx *cli.Context) error {
		fmt.Println(f())
		return nil
	}
}

func run(f func() error) interface{} {
	return func(ctx *cli.Context) error {
		return f()
	}
}
