package mylibrary

import (
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary/controller"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type (
	// Router to server
	Router struct {
		dig.In
		BookCntrl controller.BookCntrl
	}
)

var _ echokit.Router = (*Router)(nil)

// SetRoute to echo server
func (r *Router) SetRoute(e echokit.Server) {
	group := e.Group("/mylibrary")
	echokit.SetRoute(group,
		&r.BookCntrl,
	)
}
