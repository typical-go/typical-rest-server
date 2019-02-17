package app

import (
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
	crud.RegisterTo(entity, s.Echo)
}
