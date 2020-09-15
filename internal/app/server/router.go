package server

import (
	"github.com/typical-go/typical-rest-server/internal/app/server/controller"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// Router to server
	Router struct {
		dig.In
		BookCntrl controller.BookCntrl
	}
)

var _ typrest.Router = (*Router)(nil)

// SetRoute to echo server
func (r *Router) SetRoute(e typrest.Server) {
	typrest.SetRoute(e,
		&r.BookCntrl,
	)
}
