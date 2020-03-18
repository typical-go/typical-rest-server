package typredis_test

import (
	"testing"

	"github.com/typical-go/typical-go/pkg/typdocker"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
)

func TestRedis(t *testing.T) {
	t.Run("SHOULD implement provider", func(t *testing.T) {
		var _ typapp.Provider = typredis.New()
	})
	t.Run("SHOULD implement destroyer", func(t *testing.T) {
		var _ typapp.Destroyer = typredis.New()
	})
	t.Run("SHOULD implement preparer", func(t *testing.T) {
		var _ typapp.Preparer = typredis.New()
	})
	t.Run("SHOULD implement composer", func(t *testing.T) {
		var _ typdocker.Composer = typredis.New()
	})
}
