package typical

import (
	"context"
	"log"
	"time"

	"github.com/typical-go/typical-rest-server/app"
)

func startApplication(s *app.Server) error {
	log.Println("Start Application")
	return s.Serve()
}

func gracefulShutdown(s *app.Server) (err error) {
	log.Println("Graceful Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Shutdown(ctx)

	return
}
