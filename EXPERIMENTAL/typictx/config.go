package typictx

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
)

const (
	defaultDotEnv = ".env"
)

// ConfigFields return config list
func ConfigFields(ctx *Context) (fields []typiobj.ConfigField) {
	if configurer, ok := ctx.AppModule.(typiobj.Configurer); ok {
		fields = append(fields, configurer.Configure().ConfigFields()...)
	}
	for _, module := range ctx.Modules {
		if configurer, ok := module.(typiobj.Configurer); ok {
			fields = append(fields, configurer.Configure().ConfigFields()...)
		}
	}
	return
}

// GenerateEnvfile to generate .env file if not exist
func GenerateEnvfile(ctx *Context) (err error) {
	if _, err = os.Stat(defaultDotEnv); !os.IsNotExist(err) {
		return
	}
	log.Infof("Generate new project environment at '%s'", defaultDotEnv)
	var file *os.File
	if file, err = os.Create(defaultDotEnv); err != nil {
		return
	}
	defer file.Close()
	for _, field := range ConfigFields(ctx) {
		s := fmt.Sprintf("%s=%s\n", field.Name, field.Default)
		file.WriteString(s)
	}
	return
}
