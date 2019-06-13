package main

import (
	"github.com/typical-go/typical-rest-server/app/server"
	"github.com/typical-go/typical-rest-server/provider"
)

func main() {
	provider.Container().Invoke(func(s *server.Server) error {
		return s.Serve()
	})
}
