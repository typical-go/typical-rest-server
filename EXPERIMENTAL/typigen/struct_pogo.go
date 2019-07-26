package typigen

import "reflect"

type StructPogo struct {
	Name   string
	Fields []reflect.StructField
}

func (s StructPogo) String() string {
	var builder Builder
	builder.Printlnf("type %s struct{", s.Name)
	for _, field := range s.Fields {
		builder.Printlnf("%s %s", field.Name, field.Type.String())
	}

	builder.Printlnf("}")

	return builder.String()
}
