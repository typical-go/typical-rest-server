package typrails

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-go/pkg/utility/coll"
	"go.uber.org/dig"
)

// Fetcher responsible to fetch entity
type Fetcher struct {
	dig.In
	*sql.DB
}

// Fetch entity from database based on table name
func (f *Fetcher) Fetch(tableName string) (e *Entity, err error) {
	var infos []InfoSchema
	if infos, err = f.infoSchema(tableName); err != nil {
		return
	}
	var fields []Field
	if fields, err = f.standardFields(infos); err != nil {
		return
	}
	forms := f.formFields(infos)
	e = &Entity{
		Name:           "music",
		Table:          "musics",
		Type:           "Music",
		Cache:          "MUSIC",
		ProjectPackage: "github.com/typical-go/typical-rest-server",
		Fields:         append(fields, forms...),
		Forms:          forms,
	}
	return
}

func (f *Fetcher) standardFields(infos []InfoSchema) (stdFields []Field, err error) {
	stdFields = []Field{
		{
			Name:      "ID",
			Column:    "id",
			Type:      "int64",
			Udt:       "int4",
			StructTag: "`json:\"id\"`",
		},
		{
			Name:      "UpdatedAt",
			Column:    "updated_at",
			Type:      "time.Time",
			Udt:       "timestamp",
			StructTag: "`json:\"updated_at\"`",
		},
		{
			Name:      "CreatedAt",
			Column:    "created_at",
			Type:      "time.Time",
			Udt:       "timestamp",
			StructTag: "`json:\"created_at\"`",
		},
	}
	var errs coll.Errors
field_loop:
	for _, field := range stdFields {
		for _, info := range infos {
			if info.ColumnName == field.Column && info.DataType == field.Udt {
				continue field_loop
			}
		}
		errs.Append(fmt.Errorf("\"%s\" with underlying data type \"%s\" is missing",
			field.Column, field.Udt))
	}
	err = errs.Unwrap()
	return
}

//
func (f *Fetcher) formFields(infos []InfoSchema) (forms []Field) {
	return
}

func (f *Fetcher) infoSchema(tableName string) (infos []InfoSchema, err error) {
	columns := sq.Select("column_name", "udt_name").
		From("information_schema.COLUMNS").
		Where(sq.Eq{"table_name": tableName})
	var rows *sql.Rows
	if rows, err = columns.RunWith(f.DB).Query(); err != nil {
		return
	}
	for rows.Next() {
		var info InfoSchema
		if err = rows.Scan(&info.ColumnName, &info.DataType); err != nil {
			return
		}
		infos = append(infos, info)
	}
	return
}

type InfoSchema struct {
	ColumnName string
	DataType   string
}
