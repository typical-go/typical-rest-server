package typcfg

import (
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typgen"
)

type (
	usageTmplData struct {
		typgen.Signature
		ProjectName string
		Fields      []*Field
		EnvSnippet  string
	}
)

const usageTmpl = `# {{.ProjectName}}

<!-- {{.Signature}} -->

## Configuration List
| Field Name | Default | Required | 
|---|---|:---:|{{range $f := .Fields}}
| {{$f.Key}} | {{$f.Default}} | {{if $f.Required}}Yes{{end}} |{{end}}

## DotEnv example
{{.EnvSnippet}}
`

// GenerateDoc generate usage document
func GenerateDoc(target string, c *Context) error {
	fields := fields(c)
	c.Infof("Generate '%s'\n", target)
	return tmplkit.WriteFile(
		target,
		usageTmpl,
		usageTmplData{
			Signature:   typgen.Signature{TagName: "@envconfig"},
			ProjectName: c.Descriptor.ProjectName,
			Fields:      fields,
			EnvSnippet:  envSnippet(fields),
		},
	)
}

func fields(c *Context) []*Field {
	var fields []*Field
	for _, config := range c.Configs {
		fields = append(fields, config.Fields...)
	}
	return fields
}

func envSnippet(fields []*Field) string {
	var env strings.Builder
	fmt.Fprintln(&env, "```")
	for _, field := range fields {
		fmt.Fprintf(&env, "%s=%s\n", field.Key, field.Default)
	}
	fmt.Fprintln(&env, "```")
	return env.String()
}
