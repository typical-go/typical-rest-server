package project

import (
	"bytes"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical"
)

const configFormat = `This application is configured via the environment. The following environment
variables can be used:

KEY	TYPE	DEFAULT	REQUIRED	DESCRIPTION
{{range .}}{{usage_key .}}	{{usage_type .}}	{{usage_default .}}	{{usage_required .}}	{{usage_description .}}
{{end}}`

// ConfigDetail return config detail string
func ConfigDetail() string {
	buf := new(bytes.Buffer)
	envconfig.Usagef(typical.Context.ConfigPrefix, &typical.AllConfig{}, buf, configFormat)
	return buf.String()
}
