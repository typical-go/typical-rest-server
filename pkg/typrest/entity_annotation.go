package typrest

import (
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-go/pkg/typast"
)

type (
	// EntityAnnotation ...
	EntityAnnotation struct {
		TagName string // By default is @entity
	}
	// Entity ...
	Entity struct {
		Repo    string
		Table   string
		Dialect string
		DBCtor  string
		Target  string
		Pkg     string
	}
)

//
// EntityAnnotation
//

var _ typast.Annotator = (*EntityAnnotation)(nil)

// Annotate Envconfig to prepare dependency-injection and env-file
func (m *EntityAnnotation) Annotate(c *typast.Context) error {
	var entities []*Entity

	annots, _ := typast.FindAnnot(c, m.getTagName(), typast.EqualStruct)
	for _, a := range annots {
		entities = append(entities, CreateEntity(a))
	}
	return nil
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
func CreateEntity(a *typast.Annot2) *Entity {
	name := a.GetName()
	tagParam := a.TagParam

	repo := tagParam.Get("repo")
	if repo == "" {
		repo = name + "Repo"
	}

	table := a.TagParam.Get("table")
	if table == "" {
		table = strings.ToLower(table) + "s"
	}

	dialect := a.TagParam.Get("dialect")

	dbCtor := a.TagParam.Get("db_ctor")

	target := a.TagParam.Get("target")
	if target == "" {
		target = filepath.Dir(a.Path)
	}
	pkg := filepath.Base(target)

	return &Entity{
		Repo:    repo,
		Table:   table,
		Dialect: dialect,
		DBCtor:  dbCtor,
		Target:  target,
		Pkg:     pkg,
	}
}
