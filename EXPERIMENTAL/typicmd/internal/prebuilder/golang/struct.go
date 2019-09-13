package golang

import (
	"io"
	"reflect"
	"strings"
)

// Struct is plain old go object for struct
type Struct struct {
	Name   string
	Fields []reflect.StructField
}

func (s Struct) Write(w io.Writer) {
	writelnf(w, "type %s struct{", s.Name)
	for _, field := range s.Fields {
		writelnf(w, "%s %s", field.Name, field.Type.String())
	}
	writeln(w, "}")
}

func (s Struct) String() string {
	var builder strings.Builder
	s.Write(&builder)
	return builder.String()
}
