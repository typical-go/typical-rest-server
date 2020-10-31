package echokit_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestResponseWriter(t *testing.T) {
	rw := echokit.NewResponseWriter()
	rw.Write([]byte("some-text"))
	rw.WriteHeader(200)
	rw.Header().Add("key", "value")

	require.Equal(t, http.Header{"Key": []string{"value"}}, rw.Header())
	require.Equal(t, "some-text", string(rw.Bytes))
	require.Equal(t, 200, rw.StatusCode)

	rw2 := echokit.NewResponseWriter()
	rw.CopyTo(rw2)
	require.Equal(t, rw, rw2)
}
