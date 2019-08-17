package typictx

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// ConfigAccessor to access the config
type ConfigAccessor interface {
	GetConfigPrefix() string
	GetConfigSpec() interface{}
	GetName() string
	GetKey() string
}

// Config represent the configuration
type Config struct {
	Prefix string
	Spec   interface{}
}

// GetConfigPrefix return the config prefix
func (c *Config) GetConfigPrefix() string {
	return c.Prefix
}

// GetConfigSpec return the config specification
func (c *Config) GetConfigSpec() interface{} {
	return c.Spec
}

// GetKey return key/field name in main config struct
func (c *Config) GetKey() string {
	if c.Prefix == "" {
		return ""
	}
	return strcase.ToCamel(strings.ToLower(c.Prefix))
}
