package typcfg

import (
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/typannot"
)

type (
	// Context of config
	Context struct {
		*typannot.Context
		Configs []*AppCfg
	}
	// AppCfg model
	AppCfg struct {
		CtorName string
		Prefix   string
		SpecType string
		Name     string
		Fields   []*Field
	}
	// Field model
	Field struct {
		Key      string
		Default  string
		Required bool
	}
)

// CreateContext create context instance
func CreateContext(c *typannot.Context, tagname string) *Context {
	var configs []*AppCfg
	for _, annot := range c.FindAnnotByStruct(tagname) {
		prefix := getPrefix(annot)
		var fields []*Field
		for _, field := range annot.Type.(*typannot.StructType).Fields {
			fields = append(fields, CreateField(prefix, field))
		}

		configs = append(configs, &AppCfg{
			CtorName: getCtorName(annot),
			Name:     annot.Name,
			Prefix:   prefix,
			SpecType: fmt.Sprintf("%s.%s", annot.Package, annot.Name),
			Fields:   fields,
		})
	}
	return &Context{Context: c, Configs: configs}
}

// CreateField create new instance of field
func CreateField(prefix string, field *typannot.Field) *Field {
	// NOTE: mimic kelseyhightower/envconfig struct tags

	name := field.Get("envconfig")
	if name == "" {
		name = strings.ToUpper(field.Name)
	}

	return &Field{
		Key:      fmt.Sprintf("%s_%s", prefix, name),
		Default:  field.Get("default"),
		Required: field.Get("required") == "true",
	}
}

func getCtorName(annot *typannot.Annot) string {
	return annot.TagParam.Get("ctor_name")
}

func getPrefix(annot *typannot.Annot) string {
	prefix := annot.TagParam.Get("prefix")
	if prefix == "" {
		prefix = strings.ToUpper(annot.Name)
	}
	return prefix
}
