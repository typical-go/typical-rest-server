package main

import (
	_ "github.com/lib/pq"
	"github.com/tiket/TIX-SESSION-GO/app"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typical.Container().Invoke(func(s *app.Server) error {
		return s.Serve()
	})
}
