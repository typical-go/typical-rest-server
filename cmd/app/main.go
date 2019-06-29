package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typical.Context.Container.Invoke(func(s *app.Server) error {
		fmt.Println("meh")
		return s.Serve()
	})
}
