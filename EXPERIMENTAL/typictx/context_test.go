package typictx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

func TestContext_AllModule(t *testing.T) {
	ctx := typictx.Context{
		AppModule: dummyApp{},
		Modules: []interface{}{
			struct{}{},
			struct{}{},
			struct{}{},
		},
	}
	require.Equal(t, 4, len(ctx.AllModule()))
	require.Equal(t, ctx.AppModule, ctx.AllModule()[0])
	for i, module := range ctx.Modules {
		require.Equal(t, module, ctx.AllModule()[i+1])
	}

}

func TestContext_Validate(t *testing.T) {
	testcases := []struct {
		context typictx.Context
		errMsg  string
	}{
		{
			typictx.Context{
				AppModule: dummyApp{},
				Name:      "some-name",
				Root:      "some-root",
				Release:   typictx.Release{Targets: []string{"linux/amd64"}},
			},
			"",
		},
		{
			typictx.Context{
				AppModule: invalidDummyApp{},
				Name:      "some-name",
				Root:      "some-root",
				Release:   typictx.Release{Targets: []string{"linux/amd64"}},
			},
			"Invalid Context: Application must implement Runner",
		},
		{
			typictx.Context{
				AppModule: dummyApp{},
				Root:      "some-root",
				Release:   typictx.Release{Targets: []string{"linux/amd64"}},
			},
			"Invalid Context: Name can't not empty",
		},
		{
			typictx.Context{
				AppModule: dummyApp{},
				Name:      "some-name",
				Release:   typictx.Release{Targets: []string{"linux/amd64"}},
			},
			"Invalid Context: Root can't not empty",
		},
		{
			typictx.Context{
				AppModule: dummyApp{},
				Name:      "some-name",
				Root:      "some-root",
			},
			"Release: Missing 'Targets'",
		},
		{
			typictx.Context{
				AppModule: dummyApp{},
				Name:      "some-name",
				Root:      "some-root",
				Release:   typictx.Release{Targets: []string{"invalid"}},
			},
			"Release: Missing '/' in target 'invalid'",
		},
	}
	for i, tt := range testcases {
		msg := fmt.Sprintf("Failed in case-%d", i)
		err := tt.context.Validate()
		if tt.errMsg == "" {
			require.NoError(t, err, msg)
		} else {
			require.EqualError(t, err, tt.errMsg, msg)
		}
	}
}

type invalidDummyApp struct {
}

type dummyApp struct {
}

func (dummyApp) Run() (runFn interface{}) {
	return
}
