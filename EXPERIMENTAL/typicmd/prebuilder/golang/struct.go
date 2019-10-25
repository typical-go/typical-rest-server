package golang

import (
	"io"
	"strings"
)

// Struct is plain old go object for struct
type Struct struct {
	Name        string
	Fields      []Field
	Description string
}

// Field is struct field
type Field struct {
	Name string
	Type string
}

// AddField to add field to struct
func (s *Struct) AddField(name, typ string) {
	s.Fields = append(s.Fields, Field{Name: name, Type: typ})
}

func (s Struct) Write(w io.Writer) {
	writelnf(w, "// %s %s", s.Name, s.Description)
	writelnf(w, "type %s struct{", s.Name)
	for _, field := range s.Fields {
		writelnf(w, "%s %s", field.Name, field.Type)
	}
	writeln(w, "}")
}

func (s Struct) String() string {
	var builder strings.Builder
	s.Write(&builder)
	return builder.String()
}
