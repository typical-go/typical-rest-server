package typienv

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const (
	defaultDotEnv = ".env"
	configKey     = "CONFIG"
	envTemplate   = `{{range .}}{{usage_key .}}={{usage_default .}}
{{end}}`
)

// LoadEnv load environment from .env file
func LoadEnv() (err error) {
	configSource := os.Getenv(configKey)
	var configs []string
	var envMap map[string]string
	if configSource == "" {
		envMap, _ = godotenv.Read()
	} else {
		configs = strings.Split(configSource, ",")
		envMap, err = godotenv.Read(configs...)
		if err != nil {
			return
		}
	}

	var builder strings.Builder
	if len(envMap) > 0 {
		builder.WriteString(fmt.Sprintf("Read the environment %s\n", configSource))
		for key, value := range envMap {
			err = os.Setenv(key, value)
			builder.WriteString(" +" + key)
		}
		log.Info(builder.String())
	}

	return
}

// WriteEnvIfNotExist will write .env file if not exist
func WriteEnvIfNotExist(ctx typictx.Context) (err error) {
	_, err = os.Stat(defaultDotEnv)
	if !os.IsNotExist(err) {
		return
	}
	log.Infof("Generate new project environment at '%s'", defaultDotEnv)

	buf, err := os.Create(defaultDotEnv)
	if err != nil {
		return
	}
	defer buf.Close()

	for _, cfg := range ctx.Configs {
		envconfig.Usagef(cfg.Prefix, cfg.Spec, buf, envTemplate)
	}

	return
}
