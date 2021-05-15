package typdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typgen"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
)

type (
	// DBRepoAnnot ...
	DBRepoAnnot struct {
		TagName string // By default is @dbrepo
	}
	// EntityTmplData ...
	EntityTmplData struct {
		typgen.Signature
		Name       string
		Table      string
		Dialect    string
		CtorDB     string
		Pkg        string
		SourcePkg  string
		Dest       string
		Fields     []*Field
		Imports    map[string]string
		PrimaryKey *Field
	}
	// Field repo
	Field struct {
		Name         string
		Type         string
		Column       string
		PrimaryKey   bool
		DefaultValue string
		SkipUpdate   bool
	}
	fieldOptions []string
)

const (
	pkOpt       = "pk"
	nowOpt      = "now"
	noUpdateOpt = "no_update"
	parentDest  = "internal/generated/dbrepo"
)

//
// DBRepoAnnot
//

var _ typgen.Processor = (*DBRepoAnnot)(nil)

// Annotate Envconfig to prepare dependency-injection and env-file
func (m *DBRepoAnnot) Process(c *typgo.Context, directive typgen.Directives) error {
	return m.Annotation().Process(c, directive)
}

func (m *DBRepoAnnot) Annotation() *typgen.Annotation {
	return &typgen.Annotation{
		Filter: typgen.Filters{
			&typgen.TagNameFilter{m.getTagName()},
			&typgen.StructFilter{},
			&typgen.PublicFilter{},
		},
		ProcessFn: m.process,
	}
}

func (m *DBRepoAnnot) process(c *typgo.Context, directive typgen.Directives) error {
	os.RemoveAll(parentDest)
	for _, directive := range directive {
		ent, err := m.createEntity(directive)
		if err != nil {
			return err
		}
		if err := m.processEnt(c, ent); err != nil {
			c.Infof("WARN: Failed process @dbrepo at '%s': %s\n", directive.GetName(), err.Error())
		}
		m.mock(c, directive, ent)
	}
	return nil
}

func (m *DBRepoAnnot) mock(c *typgo.Context, a *typgen.Directive, ent *EntityTmplData) error {
	destPkg := filepath.Base(ent.Dest)
	dest := fmt.Sprintf("%s/%s_repo_mock.go", ent.Dest, strings.ToLower(ent.Name))
	pkg := fmt.Sprintf("%s/%s", typgo.ProjectPkg, ent.Dest)
	name := ent.Name + "Repo"

	return typmock.MockGen(c, destPkg, dest, pkg, name)
}

func (m *DBRepoAnnot) processEnt(c *typgo.Context, ent *EntityTmplData) error {
	tmpl, err := getTemplate(ent.Dialect)
	if err != nil {
		return err
	}

	os.MkdirAll(ent.Dest, 0777)
	path := fmt.Sprintf("%s/%s_repo.go", ent.Dest, strings.ToLower(ent.Name))
	c.Infof("Generate repository: %s\n", path)
	if err := tmplkit.WriteFile(path, tmpl, ent); err != nil {
		return err
	}
	typgo.GoImports(c, path)
	return nil
}

func getTemplate(dialect string) (string, error) {
	switch strings.ToLower(dialect) {
	case "postgres":
		return postgresTmpl, nil
	case "mysql":
		return mysqlTmpl, nil
	}
	return "", fmt.Errorf("unknown dialect: %s", dialect)
}

func (m *DBRepoAnnot) getTagName() string {
	if m.TagName == "" {
		m.TagName = "@dbrepo"
	}
	return m.TagName
}

//
// Entity
//

// CreateEntity create entity
func (m *DBRepoAnnot) createEntity(directive *typgen.Directive) (*EntityTmplData, error) {
	name := directive.GetName()
	table := directive.TagParam.Get("table")

	if table == "" {
		table = strings.ToLower(name) + "s"
	}

	dialect := directive.TagParam.Get("dialect")

	ctorDB := directive.TagParam.Get("ctor_db")
	if ctorDB != "" {
		ctorDB = fmt.Sprintf("`name:\"%s\"`", ctorDB)
	}

	dest := m.GetDest(directive.Path)
	pkg := filepath.Base(dest)
	sourcePkg := filepath.Base(filepath.Dir(directive.Path))
	fields, primaryKey := m.createFields(directive)

	imports := map[string]string{
		"context":                         "",
		"database/sql":                    "",
		"fmt":                             "",
		"time":                            "",
		"github.com/Masterminds/squirrel": "sq",
		"github.com/typical-go/typical-rest-server/pkg/sqkit":      "",
		"github.com/typical-go/typical-rest-server/pkg/dbtxn":      "",
		"github.com/typical-go/typical-rest-server/pkg/reflectkit": "",
		"github.com/typical-go/typical-go/pkg/typapp":              "",
		"go.uber.org/dig": "",
		typgo.ProjectPkg + "/" + filepath.Dir(directive.File.Path): "",
	}

	return &EntityTmplData{
		Signature:  typgen.Signature{TagName: m.getTagName()},
		Name:       name,
		Table:      table,
		Dialect:    dialect,
		CtorDB:     ctorDB,
		Pkg:        pkg,
		SourcePkg:  sourcePkg,
		Dest:       dest,
		Fields:     fields,
		PrimaryKey: primaryKey,
		Imports:    imports,
	}, nil
}

// GetDest get destination
func (*DBRepoAnnot) GetDest(file string) string {
	source := filepath.Dir(file)
	source = strings.TrimPrefix(source, "internal/")

	if strings.HasSuffix(source, "entity") {
		return parentDest
	}
	if strings.HasSuffix(source, "/") || strings.HasSuffix(source, "_") {
		return fmt.Sprintf("%s/%srepo", parentDest, source)
	}
	return fmt.Sprintf("%s/%s_repo", parentDest, source)
}

func (m *DBRepoAnnot) createFields(directive *typgen.Directive) (fields []*Field, primaryKey *Field) {
	structDecl := directive.Decl.Type.(*typgen.StructDecl)
	for _, f := range structDecl.Fields {
		name := f.Names[0]
		column := f.StructTag.Get("column")
		if column == "" {
			column = strings.ToLower(name)
		}
		var opts fieldOptions = strings.Split(f.StructTag.Get("option"), ",")
		field := &Field{
			Name:         name,
			Type:         f.Type,
			Column:       column,
			PrimaryKey:   opts.primaryKey(),
			DefaultValue: opts.defaultValue(),
			SkipUpdate:   opts.skipUpdate(),
		}
		fields = append(fields, field)
		if field.PrimaryKey {
			primaryKey = field
		}
	}
	return
}

//
// FieldOption
//

func (o fieldOptions) primaryKey() bool {
	for _, opt := range o {
		if strings.EqualFold(opt, pkOpt) {
			return true
		}
	}
	return false
}

func (o fieldOptions) defaultValue() string {
	for _, opt := range o {
		switch strings.ToLower(opt) {
		case "now":
			return "time.Now()"
		}
	}
	return ""
}

func (o fieldOptions) skipUpdate() bool {
	for _, opt := range o {
		if strings.EqualFold(opt, noUpdateOpt) {
			return true
		}
	}
	return false
}
