package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typicli"
	"github.com/typical-go/typical-go/pkg/typimodule"
	"github.com/typical-go/typical-rest-server/app"
)

func TestModule(t *testing.T) {
	m := app.Module()
	require.True(t, typimodule.IsProvider(m))
	require.True(t, typimodule.IsConfigurer(m))
	require.True(t, typimodule.IsPreparer(m))
	require.True(t, typimodule.IsRunner(m))
	require.True(t, typicli.IsAppCommander(m))
}
