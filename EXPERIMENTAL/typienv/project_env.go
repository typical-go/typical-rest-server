package typienv

import (
	"bufio"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const (
	envFile     = ".env" // TODO: .env shoud be store in typienv
	envTemplate = `{{range .}}export {{usage_key .}}={{usage_default .}}
{{end}}`
)

func ExportProjectEnv() (err error) {
	file, err := os.Open(envFile)
	if err != nil {
		return
	}
	defer file.Close()

	builder := strings.Builder{}
	builder.WriteString("Export the environment from " + envFile + ": ")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "export") {
			args := strings.TrimSpace(text[len("export"):])
			pair := strings.Split(args, "=")
			if len(pair) > 1 {
				os.Setenv(pair[0], pair[1])
				builder.WriteString("+" + pair[0] + " ")
			}
		}
	}

	log.Info(builder.String())
	return
}

func GenerateAppEnvIfNotExist(ctx typictx.Context) (err error) {
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
