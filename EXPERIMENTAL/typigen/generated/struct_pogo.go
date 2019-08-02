package generated

import (
	"io"
	"reflect"
	"strings"
)

// StructPogo is plain old go object for struct
type StructPogo struct {
	Name   string
	Fields []reflect.StructField
}

func (s StructPogo) Write(w io.Writer) {
	writelnf(w, "type %s struct{", s.Name)
	for _, field := range s.Fields {
		writelnf(w, "%s %s", field.Name, field.Type.String())
	}
	writeln(w, "}")
}

func (s StructPogo) String() string {
	var builder strings.Builder
	s.Write(&builder)
	return builder.String()
}
