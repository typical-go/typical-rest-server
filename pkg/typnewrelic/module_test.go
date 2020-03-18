package typnewrelic_test

import (
	"testing"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typnewrelic"
)

func TestModule(t *testing.T) {
	t.Run("SHOULD implement provider", func(t *testing.T) {
		var _ typapp.Provider = typnewrelic.New()
	})
	t.Run("SHOULD implement configurer", func(t *testing.T) {
		var _ typcfg.Configurer = typnewrelic.New()
	})

}
