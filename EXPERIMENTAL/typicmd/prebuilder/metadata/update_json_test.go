package metadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
)

func TestUpdateJSON(t *testing.T) {
	name := "example-for-test"
	defer metadata.CleanAll()
	t.Run("GIVEN no file", func(t *testing.T) {
		updated, err := metadata.Update(name, map[string]string{
			"hello": "world",
		})
		require.NoError(t, err)
		require.True(t, updated)
	})
	t.Run("GIVEN existing file", func(t *testing.T) {
		t.Run("WHEN update with same data", func(t *testing.T) {
			updated, err := metadata.Update(name, map[string]string{
				"hello": "world",
			})
			require.NoError(t, err)
			require.False(t, updated)
		})
		t.Run("WHEN update withdifferent data", func(t *testing.T) {
			updated, err := metadata.Update(name, map[string]string{
				"hello": "again",
			})
			require.NoError(t, err)
			require.True(t, updated)
		})
	})
}
