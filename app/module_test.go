package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/typical-go/typical-rest-server/app"
)

func TestModule(t *testing.T) {
	m := app.Module()
	require.True(t, typiobj.IsProvider(m))
	require.True(t, typiobj.IsConfigurer(m))
	require.True(t, typiobj.IsPreparer(m))
}
