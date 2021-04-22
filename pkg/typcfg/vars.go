package typcfg

const (
	defaultCfgTarget = "internal/generated/envcfg/envcfg.go"
)

const defaultCfgTemplate = `package {{.Package}}

/* {{.Signature}} */

import ({{range $import, $alias := .Imports}}
	{{$alias}} "{{$import}}"{{end}}
)

func init() { {{if .Configs}}{{range $c := .Configs}}
	typapp.Provide("{{$c.Ctor}}",{{$c.FnName}}){{end}}{{end}}
}
{{range $c := .Configs}}
// {{$c.FnName}} load env to new instance of {{$c.Name}}
func {{$c.FnName}}() (*{{$c.SpecType}}, error) {
	var cfg {{$c.SpecType}}
	prefix := "{{$c.Prefix}}"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}{{end}}
`
