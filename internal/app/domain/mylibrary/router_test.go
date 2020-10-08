package mylibrary_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestRoute(t *testing.T) {
	e := echo.New()
	typrest.SetRoute(e, &mylibrary.Router{})
	require.Equal(t, []string{
		"/mylibrary/books\tGET,POST",
		"/mylibrary/books/:id\tDELETE,GET,PATCH,PUT",
	}, typrest.DumpEcho(e))
}
