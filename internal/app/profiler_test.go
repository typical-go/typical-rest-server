package app_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestProfiler(t *testing.T) {
	e := echo.New()
	app.SetProfiler(e, app.HealthCheck{})

	require.Equal(t, []string{
		"/application/health\tGET,HEAD",
		"/debug/*\tGET",
		"/debug/*/*\tGET",
	}, echokit.DumpEcho(e))
}
