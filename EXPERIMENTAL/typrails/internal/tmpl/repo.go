package tmpl

// Repo template
const Repo = `package repository

import (
	"context"
	"time"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

// {{.Type}} represented  {{.Name}} entity
type {{.Type}} struct {
	{{range $field := .Fields}}{{$field.Name}} {{$field.Type}} {{$field.StructTag}}
	{{end}}
}

// {{.Type}}Repo to handle {{.Table}} entity
type {{.Type}}Repo interface {
	FindOne(context.Context, int64) (*{{.Type}}, error)
	Find(context.Context) ([]*{{.Type}}, error)
	Create(context.Context, {{.Type}}) (lastInsertID int64, err error)
	Delete(context.Context, int64) error
	Update(context.Context, {{.Type}}) error
}

// New{{.Type}}Repo return new instance of {{.Type}}Repo
func New{{.Type}}Repo(impl Cached{{.Type}}RepoImpl) {{.Type}}Repo {
	return &impl
}
`
