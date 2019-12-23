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
func (r *{{.Type}}RepoImpl) Find(ctx context.Context, id int64) (e *{{.Type}}, err error) {
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
		e = new({{.Type}})
		if err = rows.Scan({{range $field := .Fields}}&e.{{$field.Name}}, {{end}}); err != nil {
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
		var e0 {{.Type}}
		if err = rows.Scan({{range $field := .Fields}}&e0.{{$field.Name}}, {{end}}); err != nil {
			return 
		}
		list = append(list, &e0)
	}
	return
}

// Insert {{.Name}}
func (r *{{.Type}}RepoImpl) Insert(ctx context.Context, e {{.Type}}) (lastInsertID int64, err error) {
	builder := sq.
		Insert("{{.Table}}").
		Columns({{range $field := .Forms}}"{{$field.Column}}",{{end}}).
		Values({{range $field := .Forms}}e.{{$field.Name}},{{end}}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	if err = builder.QueryRowContext(ctx).Scan(&e.ID); err != nil {
		return
	}
	lastInsertID = e.ID
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
func (r *{{.Type}}RepoImpl) Update(ctx context.Context, e {{.Type}}) (err error) {
	builder := sq.
		Update("{{.Table}}").
		{{range $field := .Forms}}Set("{{$field.Column}}", e.{{$field.Name}}).
		{{end}}
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": e.ID}).
		PlaceholderFormat(sq.Dollar).RunWith(typrails.TxCtx(ctx, r))
	_, err = builder.ExecContext(ctx)
	return
}
`
