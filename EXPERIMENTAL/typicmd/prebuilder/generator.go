package prebuilder

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/metadata"
	"github.com/typical-go/typical-rest-server/pkg/utility/filekit"
)

type generator interface {
	generate(target string) error
}

// Generate the go file
func Generate(name string, g generator) (updated bool, err error) {
	target := dependency + "/" + name + ".go"
	if updated, err = metadata.Update(name, g); err != nil {
		return
	}
	updated = updated || !filekit.IsExist(target)
	if updated {
		err = g.generate(target)
	}
	return
}
