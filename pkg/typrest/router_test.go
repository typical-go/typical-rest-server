package typrest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestRouter(t *testing.T) {
	var out strings.Builder
	typrest.SetRoute(nil,
		typrest.NewRouter(func(typrest.Server) { fmt.Fprintln(&out, "1") }),
		typrest.NewRouter(func(typrest.Server) { fmt.Fprintln(&out, "2") }),
	)
	require.Equal(t, "1\n2\n", out.String())
}

func TestRoutePath(t *testing.T) {
	paths := typrest.RoutePaths(typrest.NewRouter(func(s typrest.Server) {
		s.Any("/any", nil)
		s.POST("/post", nil)
		s.GET("/get", nil)
		s.PATCH("/patch", nil)
		s.DELETE("/delete", nil)
		s.PUT("/put", nil)
	}))
	require.Equal(t, []string{
		"CONNECT /any",
		"DELETE /any",
		"DELETE /delete",
		"GET /any",
		"GET /get",
		"HEAD /any",
		"OPTIONS /any",
		"PATCH /any",
		"PATCH /patch",
		"POST /any",
		"POST /post",
		"PROPFIND /any",
		"PUT /any",
		"PUT /put",
		"REPORT /any",
		"TRACE /any",
	}, paths)
}
