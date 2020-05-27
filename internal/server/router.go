package server

import (
	"github.com/typical-go/typical-rest-server/internal/server/controller"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

var _ echokit.Router = (*Router)(nil)

type (
	// Router to server
	Router struct {
		dig.In
		BookCntrl controller.BookCntrl
	}
)

// Route to echo server
func (r *Router) Route(e echokit.Server) error {
	return echokit.SetRoute(e,
		&r.BookCntrl,
	)
}
