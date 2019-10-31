package typictx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

type invalidDummyApp struct {
}

type dummyApp struct {
}

func (dummyApp) Run() (runFn interface{}) {
	return
}

func TestContext_Validate(t *testing.T) {
	testcases := []struct {
		context typictx.Context
		errMsg  string
	}{
		{
			typictx.Context{Application: invalidDummyApp{}, Name: "some-name", Root: "some-root"},
			"Invalid Context: Application must implement Runner",
		},
		{
			typictx.Context{Application: dummyApp{}, Root: "some-root"},
			"Invalid Context: Name can't not empty",
		},
		{
			typictx.Context{Application: dummyApp{}, Name: "some-name"},
			"Invalid Context: Root can't not empty",
		},
		{
			typictx.Context{Application: dummyApp{}, Name: "some-name", Root: "some-root"},
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
