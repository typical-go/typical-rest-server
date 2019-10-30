package typictx

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	defaultDotEnv = ".env"
	configKey     = "CONFIG"
	envTemplate   = `{{range .}}{{usage_key .}}={{usage_default .}}
{{end}}`
)

// CliLoadEnvFile is cli version of LoadEnvFile
func CliLoadEnvFile(ctx *cli.Context) (err error) {
	return LoadEnvFile()
}

// LoadEnvFile to load environment from .env file
// TODO: move to util
func LoadEnvFile() (err error) {
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
