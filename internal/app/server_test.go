package app_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/internal/app/controller"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestServer(t *testing.T) {
	e := echo.New()
	app.SetServer(e,
		controller.BookCntrl{},
	)

	require.Equal(t, []string{
		"/books\tGET,POST",
		"/books/:id\tDELETE,GET,HEAD,PATCH,PUT",
	}, echokit.DumpEcho(e))
}
