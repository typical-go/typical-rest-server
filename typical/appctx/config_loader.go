package appctx

// ConfigLoader provide load configuration
type ConfigLoader interface {
	LoadFunc() interface{}
	Config() interface{}
	ConfigPrefix() string
}

// ConfigDetail contain detail of configuration
type ConfigDetail struct {
	configPrefix string
	config       interface{}
}

// NewConfigDetail return new instance of ConfigDetail
func NewConfigDetail(configPrefix string, config interface{}) ConfigDetail {
	return ConfigDetail{
		configPrefix: configPrefix,
		config:       config,
	}
}

// Config return empty configuration
func (d ConfigDetail) Config() interface{} {
	return d.config
}

// ConfigPrefix is prefix of config key when get from the source
func (d ConfigDetail) ConfigPrefix() string {
	return d.configPrefix
}
