package tmpl

// RepoImpl template
const RepoImpl = `package repository

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
