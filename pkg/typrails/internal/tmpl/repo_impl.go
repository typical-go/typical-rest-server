package tmpl

// RepoImpl template
const RepoImpl = `package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/dig"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
)

// {{.Type}}RepoImpl is implementation {{.Name}} repository
type {{.Type}}RepoImpl struct {
	dig.In
	*sql.DB
}

// Find {{.Name}}
func (r *{{.Type}}RepoImpl) Find(ctx context.Context, id int64) ({{.Name}} *{{.Type}}, err error) {
	var rows *sql.Rows
	builder := sq.
		Select({{range $field := .Fields}}"{{$field.Column}}",{{end}}).
		From("{{.Table}}").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	if rows.Next() {
		{{$.Name}} = new({{.Type}})
		{{$var_name:=.Name}}if err = rows.Scan({{range $field := .Fields}}&{{$var_name}}.{{$field.Name}}, {{end}}); err != nil {
			return nil, err
		}
	}
	return
}

// List {{.Name}}
func (r *{{.Type}}RepoImpl) List(ctx context.Context) (list []*{{.Type}}, err error) {
	var rows *sql.Rows
	builder := sq.
		Select({{range $field := .Fields}}"{{$field.Column}}",{{end}}).
		From("{{.Table}}").
		PlaceholderFormat(sq.Dollar).RunWith(r)
	if rows, err = builder.QueryContext(ctx); err != nil {
		return
	}
	list = make([]*{{.Type}}, 0)
	for rows.Next() {
		var {{$.Name}}0 {{.Type}}
		{{$var_name:=.Name}}if err = rows.Scan({{range $field := .Fields}}&{{$var_name}}0.{{$field.Name}}, {{end}}); err != nil {
			return 
		}
		list = append(list, &{{$var_name}}0)
	}
	return
}

// Insert {{.Name}}
func (r *{{.Type}}RepoImpl) Insert(ctx context.Context, {{.Name}} {{.Type}}) (lastInsertID int64, err error) {
	builder := sq.
		Insert("{{.Table}}").
		Columns({{range $field := .Forms}}"{{$field.Column}}",{{end}}).
		{{$var_name:=.Name}}Values({{range $field := .Forms}}{{$var_name}}.{{$field.Name}},{{end}}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&{{.Name}}.ID); err != nil {
		return
	}
	lastInsertID = {{.Name}}.ID
	return
}

// Delete {{.Name}}
func (r *{{.Type}}RepoImpl) Delete(ctx context.Context, id int64) (err error) {
	builder := sq.
		Delete("{{.Table}}").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}

// Update {{.Name}}
func (r *{{.Type}}RepoImpl) Update(ctx context.Context, {{.Name}} {{.Type}}) (err error) {
	builder := sq.
		Update("{{.Table}}").
		{{$var_name:=.Name}}{{range $field := .Forms}}Set("{{$field.Column}}", {{$var_name}}.{{$field.Name}}).
		{{end}}
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": {{.Name}}.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
`
