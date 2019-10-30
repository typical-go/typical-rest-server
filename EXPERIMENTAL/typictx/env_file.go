package typictx

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	defaultDotEnv = ".env"

	envTemplate = `{{range .}}{{usage_key .}}={{usage_default .}}
{{end}}`
)

// PrepareEnvFile to write .env file if not exist
func PrepareEnvFile(ctx *Context) (err error) {
	if _, err = os.Stat(defaultDotEnv); !os.IsNotExist(err) {
		return
	}
	log.Infof("Generate new project environment at '%s'", defaultDotEnv)
	var file *os.File
	if file, err = os.Create(defaultDotEnv); err != nil {
		return
	}
	defer file.Close()
	for _, cfg := range ctx.Configurations() {
		envconfig.Usagef(cfg.Prefix, cfg.Spec, file, envTemplate)
	}
	return
}
