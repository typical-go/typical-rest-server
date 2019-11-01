package typictx

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	defaultDotEnv = ".env"
)

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
	for _, cfg := range Configurations(ctx) {
		for _, field := range cfg.ConfigFields() {
			s := fmt.Sprintf("%s=%s\n", field.Name, field.Default)
			file.WriteString(s)
		}
	}
	return
}
