package mylibrary

import (
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary/controller"
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
	group := e.Group("/mylibrary")
	typrest.SetRoute(group,
		&r.BookCntrl,
	)
}
