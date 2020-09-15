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
