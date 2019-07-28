package typictx

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// Config represent configuration that need for this project
type Config struct {
	Prefix      string
	Spec        interface{}
	Description string
}

// CamelPrefix return config prefix in camel case
func (c Config) CamelPrefix() string {
	return strcase.ToCamel(strings.ToLower(c.Prefix))
}
