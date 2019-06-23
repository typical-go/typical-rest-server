package extension

import (
	"fmt"

	"github.com/typical-go/typical-rest-server/typical"
	"gopkg.in/urfave/cli.v1"
)

// ActionTrigger provider common action
type ActionTrigger struct{}

// Invoke the function with DI container
func (t ActionTrigger) Invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := typical.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}

// NotImplement for not implemented function
func (t ActionTrigger) NotImplement(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

// Print trigger function and print the result
func (t ActionTrigger) Print(f func() string) interface{} {
	return func(ctx *cli.Context) error {
		fmt.Println(f())
		return nil
	}
}

// Run the function
func (t ActionTrigger) Run(f func() error) interface{} {
	return func(ctx *cli.Context) error {
		return f()
	}
}
