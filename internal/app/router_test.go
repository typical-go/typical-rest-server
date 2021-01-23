package app_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestRoute(t *testing.T) {
	e := echo.New()
	echokit.SetRoute(e, &app.Router{})
	require.Equal(t, []string{
		"/books\tGET,POST",
		"/books/:id\tDELETE,GET,HEAD,PATCH,PUT",
		"/songs\tGET,POST",
		"/songs/:id\tDELETE,GET,HEAD,PATCH,PUT",
	}, echokit.DumpEcho(e))
}
