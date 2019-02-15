package app

import (
	"fmt"

	"github.com/imantung/typical-go-server/app/controller"
	"github.com/labstack/echo"
)

type server struct {
	*echo.Echo
	bookController controller.BookController
}

func newServer(bookController controller.BookController) *server {
	s := &server{
		Echo:           echo.New(),
		bookController: bookController,
	}

	initMiddlewares(s)
	initRoutes(s)

	return s
}

func (s *server) CRUD(entity string, crud controller.CRUD) {
	s.GET(fmt.Sprintf("/%s", entity), crud.List)
	s.POST(fmt.Sprintf("/%s", entity), crud.Create)
	s.GET(fmt.Sprintf("/%s/:id", entity), crud.Get)
	s.PUT(fmt.Sprintf("/%s/:id", entity), crud.Update)
	s.DELETE(fmt.Sprintf("/%s/:id", entity), crud.Delete)
}
