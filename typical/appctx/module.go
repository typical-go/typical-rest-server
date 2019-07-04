package appctx

import (
	"gopkg.in/urfave/cli.v1"
)

// Module in typical-go applicaiton
type Module interface {
	Name() string
	ShortName() string
	ConfigPrefix() string
	Config() interface{}
	Constructors() []interface{}
	Command() cli.Command
	LoadFunc() interface{}
}

// ModuleDetail contain detail of module
type ModuleDetail struct {
	name         string
	shortName    string
	configPrefix string
	config       interface{}
}

// NewModuleDetail return new instance of ModuleDetail
func NewModuleDetail() *ModuleDetail {
	return &ModuleDetail{}
}

// SetName to set name field
func (d *ModuleDetail) SetName(name string) *ModuleDetail {
	d.name = name
	return d
}

// SetShortName to set shortName field
func (d *ModuleDetail) SetShortName(shortName string) *ModuleDetail {
	d.shortName = shortName
	return d
}

// SetConfigPrefix to set configPrefix field
func (d *ModuleDetail) SetConfigPrefix(configPrefix string) *ModuleDetail {
	d.configPrefix = configPrefix
	return d
}

// SetConfig to set config field
func (d *ModuleDetail) SetConfig(config interface{}) *ModuleDetail {
	d.config = config
	return d
}

// Name field value
func (d *ModuleDetail) Name() string {
	return d.name
}

// ShortName field value
func (d *ModuleDetail) ShortName() string {
	return d.shortName
}

// ConfigPrefix field value
func (d *ModuleDetail) ConfigPrefix() string {
	return d.configPrefix
}

// Config field value
func (d *ModuleDetail) Config() interface{} {
	return d.config
}
