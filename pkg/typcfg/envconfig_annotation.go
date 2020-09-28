package typcfg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	// EnvconfigAnnotation handle @envconfig annotation
	// e.g. `@envconfig (prefix: "PREFIX" ctor_name:"CTOR")`
	EnvconfigAnnotation struct {
		TagName  string // By default is `@envconfig`
		Template string // By default defined in defaultCfgTemplate variable
		Target   string // By default is `cmd/PROJECT_NAME/envconfig_annotated.go`
		DotEnv   string // Dotenv path. It will be generated if not empty
		UsageDoc string // Usage documentation path. It will be if not emtpy
	}
	// EnvconfigTmplData template
	EnvconfigTmplData struct {
		typast.Signature
		Package string
		Imports []string
		Configs []*Envconfig
	}
	// Context of config
	Context struct {
		*typast.Context
		Configs []*Envconfig
	}
	// Envconfig model
	Envconfig struct {
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

// Stdout standard output
var Stdout io.Writer = os.Stdout

const defaultCfgTemplate = `package {{.Package}}

/* {{.Signature}}*/

import ({{range $import := .Imports}}
	"{{$import}}"{{end}}
)

func init() { {{if .Configs}}
	typapp.AppendCtor({{range $c := .Configs}}
		&typapp.Constructor{Name: "{{$c.CtorName}}",Fn: Load{{$c.Name}}},{{end}}
	){{end}}
}
{{range $c := .Configs}}
// Load{{$c.Name}} load env to new instance of {{$c.Name}}
func Load{{$c.Name}}() (*{{$c.SpecType}}, error) {
	var cfg {{$c.SpecType}}
	prefix := "{{$c.Prefix}}"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}{{end}}
`

//
// EnvconfigAnnotation
//

var _ typast.Annotator = (*EnvconfigAnnotation)(nil)

// Annotate Envconfig to prepare dependency-injection and env-file
func (m *EnvconfigAnnotation) Annotate(c *typast.Context) error {
	context := m.Context(c)

	if err := m.generate(context); err != nil {
		return err
	}

	if m.DotEnv != "" {
		if err := GenerateAndLoadDotEnv(m.DotEnv, context); err != nil {
			return err
		}
	}

	if m.UsageDoc != "" {
		if err := GenerateUsage(m.UsageDoc, context); err != nil {
			return err
		}
	}

	return nil
}

// Context create context instance
func (m *EnvconfigAnnotation) Context(c *typast.Context) *Context {
	var configs []*Envconfig
	for _, annot := range c.FindAnnot(m.getTagName(), typast.EqualStruct) {
		configs = append(configs, createEnvconfig(annot))
	}
	return &Context{Context: c, Configs: configs}
}

func (m *EnvconfigAnnotation) generate(c *Context) error {
	target := m.getTarget(c)
	if len(c.Configs) < 1 {
		os.Remove(target)
		return nil
	}

	imports := createImports(c.Dirs)
	imports = append(imports, "github.com/kelseyhightower/envconfig")

	fmt.Fprintf(Stdout, "Generate @envconfig to %s\n", target)
	if err := tmplkit.WriteFile(target, m.getTemplate(), &EnvconfigTmplData{
		Signature: typast.Signature{
			TagName: m.getTagName(),
			Help:    "https://pkg.go.dev/github.com/typical-go/typical-rest-server/pkg/typcfg",
		},
		Package: filepath.Base(c.Destination),
		Imports: imports,
		Configs: c.Configs,
	}); err != nil {
		return err
	}
	typgo.GoImports(target)
	return nil
}

func (m *EnvconfigAnnotation) getTagName() string {
	if m.TagName == "" {
		m.TagName = "@envconfig"
	}
	return m.TagName
}

func (m *EnvconfigAnnotation) getTemplate() string {
	if m.Template == "" {
		m.Template = defaultCfgTemplate
	}
	return m.Template
}

func (m *EnvconfigAnnotation) getTarget(c *Context) string {
	if m.Target == "" {
		m.Target = fmt.Sprintf("%s/envconfig_annotated.go", c.Destination)
	}
	return m.Target
}

func createImports(dirs []string) []string {
	var imports []string
	for _, dir := range dirs {
		imports = append(imports, fmt.Sprintf("%s/%s", typgo.ProjectPkg, dir))
	}
	return imports
}

func createEnvconfig(annot *typast.Annot) *Envconfig {
	prefix := getPrefix(annot)
	structDecl := annot.Type.(*typast.StructDecl)

	name := annot.GetName()
	return &Envconfig{
		CtorName: getCtorName(annot),
		Name:     name,
		Prefix:   prefix,
		SpecType: fmt.Sprintf("%s.%s", annot.Package, name),
		Fields:   createFields(structDecl, prefix),
	}
}

func createFields(structDecl *typast.StructDecl, prefix string) []*Field {
	var fields []*Field
	for _, field := range structDecl.Fields {
		fields = append(fields, CreateField(prefix, field))
	}
	return fields
}

// CreateField create new instance of field
func CreateField(prefix string, field *typast.Field) *Field {
	// NOTE: mimic kelseyhightower/envconfig struct tags

	name := field.Get("envconfig")
	if name == "" {
		name = strings.ToUpper(field.Names[0])
	}

	return &Field{
		Key:      fmt.Sprintf("%s_%s", prefix, name),
		Default:  field.Get("default"),
		Required: field.Get("required") == "true",
	}
}

func getCtorName(annot *typast.Annot) string {
	return annot.TagParam.Get("ctor_name")
}

func getPrefix(annot *typast.Annot) string {
	prefix := annot.TagParam.Get("prefix")
	if prefix == "" {
		prefix = strings.ToUpper(annot.GetName())
	}
	return prefix
}
