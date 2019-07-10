package typicli

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
)

const (
	envFile     = ".env" // TODO: .env shoud be store in typienv
	envTemplate = `{{range .}}export {{usage_key .}}={{usage_default .}}
{{end}}`
)

func exportEnviroment() (err error) {
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

	log.Println(builder.String())
	return
}

func generateNewEnviromentIfNotExist(ctx appctx.Context) (isGenerated bool, err error) {
	_, err = os.Stat(envFile)
	if !os.IsNotExist(err) {
		isGenerated = false
		return
	}

	buf, err := os.Create(envFile)
	if err != nil {
		return
	}

	envconfig.Usagef(ctx.TypiApp.ConfigPrefix, ctx.TypiApp.Config, buf, envTemplate)

	for i := range ctx.Modules {
		module := ctx.Modules[i]
		envconfig.Usagef(module.ConfigPrefix, module.Config, buf, envTemplate)
	}

	isGenerated = true

	return
}
