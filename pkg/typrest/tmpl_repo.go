package typrest

const repoTemplate = `package repository

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

const repoImplTemplate = `package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/dig"
)

// {{.Type}}RepoImpl is implementation {{.Name}} repository
type {{.Type}}RepoImpl struct {
	dig.In
	*sql.DB
}

// Find {{.Name}}
func (r *{{.Type}}RepoImpl) Find(ctx context.Context, id int64) ({{.Name}} *{{.Type}}, err error) {
	var rows *sql.Rows
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select({{range $field := .Fields}}"{{$field.Column}}",{{end}}).
		From("{{.Table}}").
		Where(sq.Eq{"id": id})
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		{{$var_name:=.Name}}var {{$var_name}}0 {{.Type}}
		if err = rows.Scan({{range $field := .Fields}}&{{$var_name}}0.{{$field.Name}}, {{end}}); err != nil {
			return nil, err
		}
	}
	return
}

// List {{.Name}}
func (r *{{.Type}}RepoImpl) List(ctx context.Context) (list []*{{.Type}}, err error) {
	var rows *sql.Rows
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select({{range $field := .Fields}}"{{$field.Column}}",{{end}}).
		From("{{.Table}}")
	if rows, err = builder.RunWith(r.DB).QueryContext(ctx); err != nil {
		return
	}
	list = make([]*{{.Type}}, 0)
	for rows.Next() {
		{{$var_name:=.Name}}var {{$var_name}}0 {{.Type}}
		if err = rows.Scan({{range $field := .Fields}}&{{$var_name}}0.{{$field.Name}}, {{end}}); err != nil {
			return nil, err
		}
		list = append(list, &{{$var_name}}0)
	}
	return
}

// Insert {{.Name}}
func (r *{{.Type}}RepoImpl) Insert(ctx context.Context, {{.Name}} {{.Type}}) (lastInsertID int64, err error) {
	{{$var_name:=.Name}}query := sq.Insert("{{.Table}}").
		Columns({{range $field := .Forms}}"{{$field.Column}}",{{end}}).
		Values({{range $field := .Forms}}{{$var_name}}.{{$field.Name}},{{end}}).
		Suffix("RETURNING \"id\"").
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)
	if err = query.QueryRowContext(ctx).Scan(&{{.Name}}.ID); err != nil {
		return
	}
	lastInsertID = {{.Name}}.ID
	return
}

// Delete {{.Name}}
func (r *{{.Type}}RepoImpl) Delete(ctx context.Context, id int64) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("{{.Table}}").Where(sq.Eq{"id": id})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}

// Update {{.Name}}
func (r *{{.Type}}RepoImpl) Update(ctx context.Context, {{.Name}} {{.Type}}) (err error) {
	{{$var_name:=.Name}}psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("{{.Table}}").
		{{range $field := .Forms}}Set("{{$field.Column}}", {{$var_name}}.{{$field.Name}}).
		{{end}}
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": {{.Name}}.ID})
	_, err = builder.RunWith(r.DB).ExecContext(ctx)
	return
}
`

const cachedRepoImplTemplate = `package repository

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/utility/cachekit"
	"go.uber.org/dig"
)

// Cached{{.Type}}RepoImpl is cached implementation of {{.Name}} repository
type Cached{{.Type}}RepoImpl struct {
	dig.In
	{{.Type}}RepoImpl
	Redis *redis.Client
}

// Find {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) Find(ctx context.Context, id int64) ({{.Name}} *{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:FIND:%d", id)
	{{.Name}} = new({{.Type}})
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, {{.Name}}); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if {{.Name}}, err = r.{{.Type}}RepoImpl.Find(ctx, id); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, {{.Name}}, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}

// List of {{.Name}} entity
func (r *Cached{{.Type}}RepoImpl) List(ctx context.Context) (list []*{{.Type}}, err error) {
	cacheKey := fmt.Sprintf("{{.Cache}}:LIST")
	redisClient := r.Redis.WithContext(ctx)
	if err = cachekit.Get(redisClient, cacheKey, &list); err == nil {
		log.Infof("Using cache %s", cacheKey)
		return
	}
	if list, err = r.{{.Type}}RepoImpl.List(ctx); err != nil {
		return
	}
	if err2 := cachekit.Set(redisClient, cacheKey, list, 20*time.Second); err2 != nil {
		log.Fatal(err2.Error())
	}
	return
}
`
