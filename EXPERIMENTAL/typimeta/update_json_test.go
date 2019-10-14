package typimeta_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typimeta"
)

func TestUpdateJSON(t *testing.T) {
	filename := "example-for-test.json"
	defer os.Remove(filename)
	t.Run("GIVEN no file", func(t *testing.T) {
		updated, err := typimeta.UpdateJSON(filename, map[string]string{
			"hello": "world",
		})
		require.NoError(t, err)
		require.True(t, updated)
	})
	t.Run("GIVEN existing file", func(t *testing.T) {
		t.Run("WHEN update with same data", func(t *testing.T) {
			updated, err := typimeta.UpdateJSON(filename, map[string]string{
				"hello": "world",
			})
			require.NoError(t, err)
			require.False(t, updated)
		})
		t.Run("WHEN update withdifferent data", func(t *testing.T) {
			updated, err := typimeta.UpdateJSON(filename, map[string]string{
				"hello": "again",
			})
			require.NoError(t, err)
			require.True(t, updated)
		})
	})
}
