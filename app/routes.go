package app

import (
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/app/controller"
)

// Routes of API
func Routes(
	server *echo.Echo,
	bookCntlr *controller.BookController,
	appCntlr *controller.ApplicationController,
) {
	log.Info("Initiate API Routes")

	bookCntlr.Route(server)
	appCntlr.Route(server)
}
