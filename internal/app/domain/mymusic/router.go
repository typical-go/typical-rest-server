package mymusic

import (
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/controller"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// Router to server
	Router struct {
		dig.In
		SongCntrl controller.SongCntrl
	}
)

var _ typrest.Router = (*Router)(nil)

// SetRoute to echo server
func (r *Router) SetRoute(e typrest.Server) {
	group := e.Group("/mymusic")
	typrest.SetRoute(group,
		&r.SongCntrl,
	)
}
