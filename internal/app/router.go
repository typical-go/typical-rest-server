package app

import (
	"github.com/typical-go/typical-rest-server/internal/app/controller"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type (
	// Router to server
	Router struct {
		dig.In
		BookCntrl controller.BookCntrl
		SongCntrl controller.SongCntrl
	}
)

var _ echokit.Router = (*Router)(nil)

// SetRoute to echo server
func (r *Router) SetRoute(e echokit.Server) {
	echokit.SetRoute(e,
		&r.BookCntrl,
		&r.SongCntrl,
	)
}
