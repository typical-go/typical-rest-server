package typictx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func TestContext_Validate(t *testing.T) {
	testcases := []struct {
		context typictx.Context
		errMsg  string
	}{
		{
			typictx.Context{},
			"Invalid Context: Name can't not empty",
		},
		{
			typictx.Context{Name: "some-name"},
			"Invalid Context: Root can't not empty",
		},
		{
			typictx.Context{Name: "some-name", Root: "some-root"},
			"",
		},
	}
	for _, tt := range testcases {
		err := tt.context.Validate()
		if tt.errMsg == "" {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, tt.errMsg)
		}
	}
}
