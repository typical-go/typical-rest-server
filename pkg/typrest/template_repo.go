package typrest

const repoTmpl = `package repository

import (
	"context"
	"time"
)

// {{ .Name }} represented database model
type {{ .Name }} struct {
	ID        int64     
	Title     string    
	Author    string    
	UpdatedAt time.Time 
	CreatedAt time.Time 
}

// {{ .Name }}Repo to get book data from database
type {{ .Name }}Repo interface {
	Find(ctx context.Context, id int64) (*{{ .Name }} , error)
	List(ctx context.Context) ([]*{{ .Name }} , error)
	Insert(ctx context.Context, {{ .SmallCase }} {{ .Name }} ) (lastInsertID int64, err error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, {{ .SmallCase }} {{ .Name }} ) error
}
`
