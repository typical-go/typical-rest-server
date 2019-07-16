package typical

import (
	"context"
	"log"
	"time"

	"github.com/typical-go/typical-rest-server/app"
)

func startApplication(s *app.Server) error {
	log.Println("Start the application")
	return s.Serve()
}

func stopApplication(s *app.Server) (err error) {
	log.Println("Stop the application")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Shutdown(ctx)
	return
}
