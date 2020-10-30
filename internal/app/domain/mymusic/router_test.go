package mymusic_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestRoute(t *testing.T) {
	e := echo.New()
	echokit.SetRoute(e, &mymusic.Router{})
	require.Equal(t, []string{
		"/mymusic/songs\tGET,POST",
		"/mymusic/songs/:id\tDELETE,GET,HEAD,PATCH,PUT",
	}, echokit.DumpEcho(e))
}
