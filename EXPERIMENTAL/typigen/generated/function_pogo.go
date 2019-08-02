package generated

import (
	"io"
	"strings"
)

// FunctionPogo is plain old go object for function
type FunctionPogo struct {
	Name         string
	FuncParams   map[string]string
	ReturnValues []string
	FuncBody     string
}

func (p FunctionPogo) Write(w io.Writer) {
	writelnf(w, "func %s(%s) (%s){ ",
		p.Name,
		p.funcParamsString(),
		strings.Join(p.ReturnValues, ","),
	)
	writelnf(w, p.FuncBody)
	write(w, "}")
}

func (p FunctionPogo) String() string {
	var builder strings.Builder
	p.Write(&builder)
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
