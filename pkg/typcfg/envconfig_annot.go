package typcfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typgen"
	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	// EnvconfigAnnot handle @envconfig annotation
	// e.g. `@envconfig (prefix: "PREFIX" ctor:"CTOR")`
	EnvconfigAnnot struct {
		TagName   string // By default is `@envconfig`
		Template  string // By default defined in defaultCfgTemplate variable
		Target    string // By default is `internal/generated/envcfg/envcfg.go`
		GenDotEnv string // Dotenv path. It will be generated if not empty
		GenDoc    string // Usage documentation path. It will be if not empty
	}
	// EnvconfigTmplData template
	EnvconfigTmplData struct {
		typgen.Signature
		Package string
		Configs []*Envconfig
		Imports map[string]string
	}
	// Context of config
	Context struct {
		*typgo.Context
		Configs []*Envconfig
		Imports map[string]string
	}
	// Envconfig model
	Envconfig struct {
		Ctor     string
		Prefix   string
		SpecType string
		Name     string
		Fields   []*Field
		FnName   string
	}
	// Field model
	Field struct {
		Key      string
		Default  string
		Required bool
	}
)

//
// EnvconfigAnnot
//

var _ typgen.Processor = (*EnvconfigAnnot)(nil)

// Annotate Envconfig to prepare dependency-injection and env-file
func (m *EnvconfigAnnot) Process(c *typgo.Context, directives typgen.Directives) error {
	return m.Annotation().Process(c, directives)
}

func (m *EnvconfigAnnot) Annotation() *typgen.Annotation {
	return &typgen.Annotation{
		Filter: typgen.Filters{
			&typgen.TagNameFilter{m.getTagName()},
			&typgen.PublicFilter{},
			&typgen.StructFilter{},
		},
		ProcessFn: m.process,
	}
}

func (m *EnvconfigAnnot) process(c *typgo.Context, directives typgen.Directives) error {
	context := m.Context(c, directives)
	target := m.getTarget(context)

	if len(context.Configs) < 1 {
		os.Remove(target)
	} else if err := m.generate(context, target); err != nil {
		return err
	}

	if m.GenDotEnv != "" {
		if err := GenerateAndLoadDotEnv(m.GenDotEnv, context); err != nil {
			return err
		}
	}

	if m.GenDoc != "" {
		if err := GenerateDoc(m.GenDoc, context); err != nil {
			return err
		}
	}
	return nil
}

// Context create context instance
func (m *EnvconfigAnnot) Context(c *typgo.Context, directive typgen.Directives) *Context {
	var configs []*Envconfig

	importAliases := typgen.NewImportAliases()
	for _, a := range directive {
		importAlias := importAliases.Append(a.Package())
		configs = append(configs, createEnvconfig(a, importAlias))
	}
	importAliases.Map["github.com/kelseyhightower/envconfig"] = ""
	importAliases.Map["github.com/typical-go/typical-go/pkg/typapp"] = ""
	importAliases.Map["fmt"] = ""

	return &Context{Context: c, Configs: configs, Imports: importAliases.Map}
}

func (m *EnvconfigAnnot) generate(c *Context, target string) error {

	dest := filepath.Dir(target)
	os.MkdirAll(dest, 0777)

	c.Infof("Generate @envconfig to %s\n", target)
	if err := tmplkit.WriteFile(target, m.getTemplate(), &EnvconfigTmplData{
		Signature: typgen.Signature{TagName: m.getTagName()},
		Package:   filepath.Base(dest),
		Imports:   c.Imports,
		Configs:   c.Configs,
	}); err != nil {
		return err
	}
	typgo.GoImports(c.Context, target)
	return nil
}

func (m *EnvconfigAnnot) getTagName() string {
	if m.TagName == "" {
		m.TagName = "@envconfig"
	}
	return m.TagName
}

func (m *EnvconfigAnnot) getTemplate() string {
	if m.Template == "" {
		m.Template = defaultCfgTemplate
	}
	return m.Template
}

func (m *EnvconfigAnnot) getTarget(c *Context) string {
	if m.Target == "" {
		m.Target = defaultCfgTarget
	}
	return m.Target
}

func createEnvconfig(a *typgen.Directive, importAlias string) *Envconfig {
	prefix := getPrefix(a)
	structDecl := a.Type.(*typgen.StructDecl)

	name := a.GetName()
	ctor := getCtorName(a)
	return &Envconfig{
		Ctor:     ctor,
		Name:     name,
		Prefix:   prefix,
		SpecType: fmt.Sprintf("%s.%s", importAlias, name),
		Fields:   createFields(structDecl, prefix),
		FnName:   fmt.Sprintf("Load%s%s", strcase.ToCamel(ctor), name),
	}
}

func createFields(structDecl *typgen.StructDecl, prefix string) []*Field {
	var fields []*Field
	for _, field := range structDecl.Fields {
		fields = append(fields, CreateField(prefix, field))
	}
	return fields
}

// CreateField create new instance of field
func CreateField(prefix string, field *typgen.Field) *Field {
	// NOTE: mimic kelseyhightower/envconfig struct tags
	name := field.Get("envconfig")
	if name == "" {
		name = strings.ToUpper(field.Names[0])
	}

	key := name
	if prefix != "" {
		key = fmt.Sprintf("%s_%s", prefix, name)
	}

	return &Field{
		Key:      key,
		Default:  field.Get("default"),
		Required: field.Get("required") == "true",
	}
}

func getCtorName(annot *typgen.Directive) string {
	return annot.TagParam.Get("ctor")
}

func getPrefix(annot *typgen.Directive) string {
	prefix := annot.TagParam.Get("prefix")
	if prefix == "" {
		return strings.ToUpper(annot.GetName())
	}
	if prefix == "-" || prefix == "_" {
		return ""
	}
	return prefix
}
