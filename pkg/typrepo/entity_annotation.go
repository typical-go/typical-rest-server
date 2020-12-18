package typrepo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-go/pkg/oskit"
	"github.com/typical-go/typical-go/pkg/tmplkit"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
)

type (
	// EntityAnnotation ...
	EntityAnnotation struct {
		TagName string // By default is @entity
	}
	// EntityTmplData ...
	EntityTmplData struct {
		typast.Signature
		Name       string
		Table      string
		Dialect    string
		CtorDB     string
		Package    string
		DestFolder string
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
)

//
// EntityAnnotation
//

var _ typast.Annotator = (*EntityAnnotation)(nil)

// Annotate Envconfig to prepare dependency-injection and env-file
func (m *EntityAnnotation) Annotate(c *typast.Context) error {
	annots, _ := typast.FindAnnot(c, m.getTagName(), typast.EqualStruct)
	for _, a := range annots {
		ent, err := m.createEntity(a)
		if err != nil {
			return err
		}
		if err := m.process(ent); err != nil {
			fmt.Fprintf(oskit.Stdout, "WARN: Failed process @entity at '%s': %s\n", a.GetName(), err.Error())
		}

		m.mock(c, ent)
	}
	return nil
}

func (m *EntityAnnotation) mock(c *typast.Context, ent *EntityTmplData) error {
	destPkg := ent.Package + "_repo_mock"
	dest := ent.DestFolder + "_mock/" + strings.ToLower(ent.Name) + "_repo.go"
	pkg := typgo.ProjectPkg + "/" + ent.DestFolder
	name := ent.Name + "Repo"

	return typmock.MockGen(c.Context, destPkg, dest, pkg, name)
}

func (m *EntityAnnotation) process(ent *EntityTmplData) error {
	tmpl, err := getTemplate(ent.Dialect)
	if err != nil {
		return err
	}

	os.MkdirAll(ent.DestFolder, 0777)
	path := fmt.Sprintf("%s/%s_repo.go", ent.DestFolder, strings.ToLower(ent.Name))
	fmt.Fprintf(oskit.Stdout, "Generate repository: %s\n", path)
	if err := tmplkit.WriteFile(path, tmpl, ent); err != nil {
		return err
	}
	typgo.GoImports(path)
	return nil
}

func getTemplate(dialect string) (string, error) {
	switch strings.ToLower(dialect) {
	case "postgres":
		return postgresTmpl, nil
	case "mysql":
		return mysqlTmpl, nil
	}
	return "", fmt.Errorf("Unknown dialect: %s", dialect)
}

func (m *EntityAnnotation) getTagName() string {
	if m.TagName == "" {
		m.TagName = "@entity"
	}
	return m.TagName
}

//
// Entity
//

// CreateEntity create entity
func (m *EntityAnnotation) createEntity(a *typast.Annot2) (*EntityTmplData, error) {
	name := a.GetName()
	table := a.TagParam.Get("table")

	if table == "" {
		table = strings.ToLower(name) + "s"
	}

	dialect := a.TagParam.Get("dialect")

	ctorDB := a.TagParam.Get("ctor_db")
	if ctorDB != "" {
		ctorDB = fmt.Sprintf("`name:\"%s\"`", ctorDB)
	}

	source := filepath.Dir(a.Path)
	pkg := filepath.Base(source)
	destFolder := fmt.Sprintf("internal/generated/entity/%s_repo", pkg)

	var fields []*Field
	var primaryKey *Field
	structDecl := a.Decl.Type.(*typast.StructDecl)
	for _, f := range structDecl.Fields {
		name := f.Names[0]
		column := f.StructTag.Get("column")
		if column == "" {
			column = strings.ToLower(name)
		}
		var opts fieldOptions
		opts = strings.Split(f.StructTag.Get("option"), ",")

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
		typgo.ProjectPkg + "/" + filepath.Dir(a.File.Path): "",
	}

	return &EntityTmplData{
		Signature:  typast.Signature{TagName: m.getTagName()},
		Name:       name,
		Table:      table,
		Dialect:    dialect,
		CtorDB:     ctorDB,
		Package:    pkg,
		DestFolder: destFolder,
		Fields:     fields,
		PrimaryKey: primaryKey,
		Imports:    imports,
	}, nil
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
