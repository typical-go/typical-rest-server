package typpostgres_test

import (
	"testing"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
)

func TestPostgres(t *testing.T) {
	t.Run("SHOULD implement commander", func(t *testing.T) {
		var _ typbuildtool.Commander = typpostgres.New()
	})
	t.Run("SHOULD implement provider", func(t *testing.T) {
		var _ typapp.Provider = typpostgres.New()
	})
	t.Run("SHOULD implement destroyer", func(t *testing.T) {
		var _ typapp.Destroyer = typpostgres.New()
	})
	t.Run("SHOULD implement preparer", func(t *testing.T) {
		var _ typapp.Preparer = typpostgres.New()
	})
	t.Run("SHOULD implement composer", func(t *testing.T) {
		var _ typdocker.Composer = typpostgres.New()
	})
	t.Run("SHOULD implement configurer", func(t *testing.T) {
		var _ typcfg.Configurer = typpostgres.New()
	})
}
