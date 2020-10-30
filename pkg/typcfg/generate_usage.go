package typcfg

import (
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typast"
)

type (
	usageTmplData struct {
		typast.Signature
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

// GenerateUsage generate usage document
func GenerateUsage(target string, c *Context) error {
	fields := fields(c)
	fmt.Fprintf(Stdout, "Generate '%s'\n", target)
	return tmplkit.WriteFile(
		target,
		usageTmpl,
		usageTmplData{
			Signature:   typast.Signature{TagName: "@envconfig"},
			ProjectName: c.BuildSys.ProjectName,
			Fields:      fields,
			EnvSnippet:  envSnippet(fields),
		},
	)
}

func fields(c *Context) []*Field {
	var fields []*Field
	for _, config := range c.Configs {
		for _, field := range config.Fields {
			fields = append(fields, field)
		}
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
