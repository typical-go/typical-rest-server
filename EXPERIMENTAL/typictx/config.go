package typictx

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type Config struct {
	Prefix      string
	Spec        interface{}
	Description string
}

func (c Config) CamelPrefix() string {
	return strcase.ToCamel(strings.ToLower(c.Prefix))
}
