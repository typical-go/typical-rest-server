package typserver_test

import (
	"testing"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

func TestModule(t *testing.T) {
	t.Run("SHOULD implement provider", func(t *testing.T) {
		var _ typapp.Provider = typserver.Module()
	})

	t.Run("SHOULD implement destroyer", func(t *testing.T) {
		var _ typapp.Destroyer = typserver.Module()
	})
}
