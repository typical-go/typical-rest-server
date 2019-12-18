package tmpl

// Repo template
const Repo = `package repository

import (
	"context"
	"time"
)

// {{.Type}} represented  {{.Name}} entity
type {{.Type}} struct {
	{{range $field := .Fields}}{{$field.Name}} {{$field.Type}} {{$field.StructTag}}
	{{end}}
}

// {{.Type}}Repo to handle {{.Name}}  entity
type {{.Type}}Repo interface {
	Find(ctx context.Context, id int64) (*{{.Type}}, error)
	List(ctx context.Context) ([]*{{.Type}}, error)
	Insert(ctx context.Context, {{.Name}} {{.Type}}) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, {{.Name}} {{.Type}}) error
}

// New{{.Type}}Repo return new instance of {{.Type}}Repo
func New{{.Type}}Repo(impl Cached{{.Type}}RepoImpl) {{.Type}}Repo {
	return &impl
}
`
