package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/app"
)

func TestModule(t *testing.T) {
	m := app.New()
	require.True(t, typcore.IsConfigurer(m))
	require.True(t, typcore.IsAppCommander(m))
}
