package typpostgres_test

import (
	"testing"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestUtility(t *testing.T) {
	t.Run("SHOULD implement Utility", func(t *testing.T) {
		var _ typbuildtool.Task = typpostgres.Utility()
	})

}
