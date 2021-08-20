package typdb

const postgresTmpl = `package {{.Pkg}}

/* {{.Signature}} */

import({{range $pkg, $alias := .Imports}}
	{{$alias}} "{{$pkg}}"{{end}}
)

var (
	// {{.Name}}TableName is table name for {{.Table}} entity
	{{.Name}}TableName = "{{.Table}}"
	// {{.Name}}Table is columns for {{.Table}} entity
	{{.Name}}Table = struct {
		{{range .Fields}}{{.Name}} string
		{{end}}
	}{
		{{range .Fields}}{{.Name}}: "{{.Column}}",
		{{end}}
	}
)

type (
	// {{.Name}}Repo to get {{.Table}} data from database
	{{.Name}}Repo interface {
		Count(context.Context, ...sqkit.SelectOption) (int64, error)
		Find(context.Context, ...sqkit.SelectOption) ([]*{{.SourcePkg}}.{{.Name}}, error)
		Insert(context.Context, *{{.SourcePkg}}.{{.Name}}) (int64, error)
		BulkInsert(context.Context, ...*{{.SourcePkg}}.{{.Name}}) (int64, error)
		Delete(context.Context, sqkit.DeleteOption) (int64, error)
		Update(context.Context, *{{.SourcePkg}}.{{.Name}}, sqkit.UpdateOption) (int64, error)
		Patch(context.Context, *{{.SourcePkg}}.{{.Name}}, sqkit.UpdateOption) (int64, error)
	}
	// {{.Name}}RepoImpl is implementation {{.Table}} repository
	{{.Name}}RepoImpl struct {
		dig.In
		*sql.DB {{.CtorDB}}
	}
)

func init() {
	typapp.Provide("",New{{.Name}}Repo)
}

// Count {{.Table}}
func (r *{{.Name}}RepoImpl) Count(ctx context.Context, opts ...sqkit.SelectOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}
	builder := sq.
		Select("count(1)").
		From({{.Name}}TableName).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	row := builder.QueryRowContext(ctx)

	var cnt int64
	if err := row.Scan(&cnt); err != nil {
		return -1, err
	}
	return cnt, nil
}

// New{{.Name}}Repo return new instance of {{.Name}}Repo
func New{{.Name}}Repo(impl {{.Name}}RepoImpl) {{.Name}}Repo {
	return &impl
}

// Find {{.Table}}
func (r *{{.Name}}RepoImpl) Find(ctx context.Context, opts ...sqkit.SelectOption) (list []*{{.SourcePkg}}.{{.Name}}, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	builder := sq.
		Select(
			{{range .Fields}}{{$.Name}}Table.{{.Name}},
			{{end}}
		).
		From({{.Name}}TableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*{{.SourcePkg}}.{{.Name}}, 0)
	for rows.Next() {
		ent := new({{.SourcePkg}}.{{.Name}})
		if err = rows.Scan({{range .Fields}}
			&ent.{{.Name}},{{end}}
		); err != nil {
			return
		}
		list = append(list, ent)
	}
	return
}

// Insert {{.Table}} and return last inserted id
func (r *{{.Name}}RepoImpl) Insert(ctx context.Context, ent *{{.SourcePkg}}.{{.Name}}) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Insert({{$.Name}}TableName).
		Columns({{range .Fields}}{{if not .PrimaryKey}}	{{$.Name}}Table.{{.Name}},{{end}}	
		{{end}}).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", {{$.Name}}Table.{{.PrimaryKey.Name}}),
		).
		PlaceholderFormat(sq.Dollar).
		Values({{range .Fields}}{{if .DefaultValue}}	{{.DefaultValue}},{{else if not .PrimaryKey}}	ent.{{.Name}},{{end}}
		{{end}})

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id {{.PrimaryKey.Type}}
	if err := scanner.Scan(&id); err != nil {
		txn.AppendError(err)
		return -1, err
	}
	return id, nil
}

// BulkInsert {{.Table}} and return affected rows
func (r *{{.Name}}RepoImpl) BulkInsert(ctx context.Context, ents ...*{{.SourcePkg}}.{{.Name}}) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Insert({{$.Name}}TableName).
		Columns({{range .Fields}}{{if not .PrimaryKey}}	{{$.Name}}Table.{{.Name}},{{end}}	
		{{end}}).
		PlaceholderFormat(sq.Dollar)
		
	
	for _, ent := range ents {
		builder = builder.Values({{range .Fields}}{{if .DefaultValue}}	{{.DefaultValue}},{{else if not .PrimaryKey}}	ent.{{.Name}},{{end}}
			{{end}})
	}

	res, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}


// Update {{.Table}}
func (r *{{.Name}}RepoImpl) Update(ctx context.Context, ent *{{.SourcePkg}}.{{.Name}}, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update({{.Name}}TableName).{{range .Fields}}{{if and (not .PrimaryKey) (not .SkipUpdate)}}
		Set({{$.Name}}Table.{{.Name}},{{if .DefaultValue}}{{.DefaultValue}}{{else}}ent.{{.Name}},{{end}}).{{end}}{{end}}
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	if opt != nil {
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}

// Patch {{.Table}}
func (r *{{.Name}}RepoImpl) Patch(ctx context.Context, ent *{{.SourcePkg}}.{{.Name}}, opt sqkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update({{.Name}}TableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	{{range .Fields}}{{if and (not .PrimaryKey) (not .SkipUpdate)}}{{if .DefaultValue}}
	builder = builder.Set({{$.Name}}Table.{{.Name}}, {{.DefaultValue}}){{else}}
	if !reflectkit.IsZero(ent.{{.Name}}) {
		builder = builder.Set({{$.Name}}Table.{{.Name}}, ent.{{.Name}})
	}{{end}}{{end}}{{end}}

	if opt != nil{
		builder = opt.CompileUpdate(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}


// Delete {{.Table}}
func (r *{{.Name}}RepoImpl) Delete(ctx context.Context, opt sqkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Delete({{.Name}}TableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	if opt != nil {
		builder = opt.CompileDelete(builder)
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.AppendError(err)
	return affectedRow, err
}
`
