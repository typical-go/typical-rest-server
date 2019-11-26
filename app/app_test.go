package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/typmodule"
	"github.com/typical-go/typical-rest-server/app"
)

func TestModule(t *testing.T) {
	m := app.Module()
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typcfg.IsConfigurer(m))
	require.True(t, typcli.IsAppCommander(m))
}
