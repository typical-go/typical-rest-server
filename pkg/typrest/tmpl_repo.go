package typrest

const repoTmpl = `package repository

import (
	"context"
	"time"
)

// {{.TypeName}} represented  {{.Name}} entity
type {{.TypeName}} struct {
	ID        int64     
	Title     string    
	Author    string    
	UpdatedAt time.Time 
	CreatedAt time.Time 
}

// {{.TypeName}}Repo to handle {{.Name}}  entity
type {{.TypeName}}Repo interface {
	Find(ctx context.Context, id int64) (*{{.TypeName}}, error)
	List(ctx context.Context) ([]*{{.TypeName}}, error)
	Insert(ctx context.Context, {{.Name}} {{.TypeName}}) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, {{.Name}} {{.TypeName}}) error
}
`
