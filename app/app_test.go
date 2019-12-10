package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-go/pkg/typobj"
	"github.com/typical-go/typical-rest-server/app"
)

func TestModule(t *testing.T) {
	m := app.Module()
	require.True(t, typobj.IsConfigurer(m))
	require.True(t, typobj.IsAppCommander(m))
}
