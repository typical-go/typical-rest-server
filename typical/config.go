package typical

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/app"
)

// LoadConfig return new instance of config
func LoadConfig() (config app.Config, err error) {
	err = envconfig.Process(Context.ConfigPrefix, &config)
	return
}
