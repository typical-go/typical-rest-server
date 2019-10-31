package runkit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown based on function
func GracefulShutdown(fn func() error) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
		fmt.Println("\n\n\nGraceful Shutdown...")
		if err := fn(); err != nil {
			fmt.Println(err.Error())
		}
	}()

}
