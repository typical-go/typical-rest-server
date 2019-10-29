package typpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
)

func TestModule(t *testing.T) {
	m := typpostgres.Module()
	require.True(t, typictx.IsConstructor(m))
	require.True(t, typictx.IsDestructor(m))
	require.True(t, typictx.IsCommandLiner(m))
	require.True(t, typictx.IsConfigurer(m))
}
