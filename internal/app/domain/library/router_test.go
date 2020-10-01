package library_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/domain/library"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestRoute(t *testing.T) {
	require.Equal(t, []string{
		"DELETE /library/books/:id",
		"GET /library/books",
		"GET /library/books/:id",
		"PATCH /library/books/:id",
		"POST /library/books",
		"PUT /library/books/:id",
	}, typrest.RoutePaths(&library.Router{}))
}
