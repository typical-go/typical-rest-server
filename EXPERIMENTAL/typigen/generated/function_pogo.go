package generated

import (
	"strings"
)

type FunctionPogo struct {
	Name         string
	FuncParams   map[string]string
	ReturnValues []string
	FuncBody     string
}

func (p FunctionPogo) String() string {
	var builder Builder
	builder.Printlnf("func %s(%s) (%s){ ",
		p.Name,
		p.funcParamsString(),
		strings.Join(p.ReturnValues, ","),
	)
	builder.Printlnf(p.FuncBody)
	builder.WriteString("}")
	return builder.String()
}

func (p FunctionPogo) funcParamsString() string {
	var builder strings.Builder
	for name := range p.FuncParams {
		builder.WriteString(name)
		builder.WriteString(" ")
		builder.WriteString(p.FuncParams[name])
		builder.WriteString(",")
	}

	return builder.String()
}
