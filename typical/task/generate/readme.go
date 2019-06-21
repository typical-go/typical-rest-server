package generate

import (
	"bytes"
	"os"
	"text/template"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-go/appx"
	"github.com/typical-go/typical-rest-server/typical"
)

var readmeData = struct {
	Context       appx.Context
	Configuration string
}{
	Context:       typical.Context,
	Configuration: configurationTable(),
}

// Readme generate readme.md
func Readme() (err error) {

	t, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return
	}

	f, err := os.Create("README.md")
	if err != nil {
		return
	}

	err = t.Execute(f, readmeData)
	if err != nil {
		return
	}
	return nil
}

func configurationTable() string {
	buf := new(bytes.Buffer)
	envconfig.Usagef(typical.Prefix, &typical.AllConfig{}, buf, `| Key | Type | Default | Request | Description |
|---|---|---|---|---|
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|
{{end}}`)
	return buf.String()
}
