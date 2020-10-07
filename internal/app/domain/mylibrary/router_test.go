package mylibrary_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

func TestRoute(t *testing.T) {
	require.Equal(t, []string{
		"DELETE /mylibrary/books/:id",
		"GET /mylibrary/books",
		"GET /mylibrary/books/:id",
		"PATCH /mylibrary/books/:id",
		"POST /mylibrary/books",
		"PUT /mylibrary/books/:id",
	}, typrest.RoutePaths(&mylibrary.Router{}))
}
