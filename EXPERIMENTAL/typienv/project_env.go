package typienv

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const (
	envFile     = ".env" // TODO: .env shoud be store in typienv
	envTemplate = `{{range .}}{{usage_key .}}={{usage_default .}}
{{end}}`
)

// LoadEnv load environment from .env file
func LoadEnv() (err error) {

	envMap, err := godotenv.Read(envFile)
	if err != nil {
		return
	}

	var builder strings.Builder
	builder.WriteString("Read the environment from '" + envFile + "':")
	for key, value := range envMap {
		err = os.Setenv(key, value)
		builder.WriteString(" +" + key)
	}
	log.Info(builder.String())

	return
}

// WriteEnvIfNotExist will write .env file if not exist
func WriteEnvIfNotExist(ctx typictx.Context) (err error) {
	_, err = os.Stat(envFile)
	if !os.IsNotExist(err) {
		return
	}
	log.Infof("Generate new project environment at '%s'", envFile)

	buf, err := os.Create(envFile)
	if err != nil {
		return
	}
	defer buf.Close()

	for _, cfg := range ctx.Configs {
		envconfig.Usagef(cfg.Prefix, cfg.Spec, buf, envTemplate)
	}

	return
}
