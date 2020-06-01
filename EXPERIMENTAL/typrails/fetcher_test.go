package typrails_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/typical-go/typical-rest-server/pkg/typrails"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestFetcher(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testcases := []struct {
		data []typrails.KeyValue
		err  error
		*typrails.Entity
	}{
		{
			data: []typrails.KeyValue{
				{"column1", "type1"},
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing; \"updated_at\" with underlying data type \"timestamp\" is missing; \"created_at\" with underlying data type \"timestamp\" is missing"),
		},
		{
			data: []typrails.KeyValue{
				{"id", "int4"},
				{"name", "varchar"},
				{"created_at", "timestamp"},
				{"updated_at", "timestamp"},
			},
			Entity: &typrails.Entity{
				Name:           "book",
				Type:           "Book",
				Table:          "books",
				Cache:          "BOOKS",
				ProjectPackage: "some-package",
				Fields: []typrails.Field{
					{Name: "ID", Type: "int64", Udt: "int4", Column: "id", StructTag: "`json:\"id\"`"},
					{Name: "Name", Type: "string", Udt: "varchar", Column: "name", StructTag: "`json:\"name\"`"},
					{Name: "CreatedAt", Type: "time.Time", Udt: "timestamp", Column: "created_at", StructTag: "`json:\"created_at\"`"},
					{Name: "UpdatedAt", Type: "time.Time", Udt: "timestamp", Column: "updated_at", StructTag: "`json:\"updated_at\"`"},
				},
				Forms: []typrails.Field{
					{Name: "Name", Type: "string", Udt: "varchar", Column: "name", StructTag: "`json:\"name\"`"},
				},
			},
		},
		{
			data: []typrails.KeyValue{
				{"id", "int"},
				{"created_at", "timestamp"},
				{"updated_at", "timestamp"},
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing"),
		},
	}
	fetcher := typrails.Fetcher{DB: db}
	query := regexp.QuoteMeta("SELECT column_name, udt_name FROM information_schema.COLUMNS WHERE table_name = $1")
	for i, tt := range testcases {
		rows := sqlmock.NewRows([]string{"column_name", "data_type"})
		for _, kv := range tt.data {
			rows.AddRow(kv.Key, kv.Value)
		}
		mock.ExpectQuery(query).WithArgs("books").WillReturnRows(rows)
		entity, err := fetcher.Fetch("some-package", "books", "book")
		require.EqualValues(t, tt.err, err, i)
		require.EqualValues(t, tt.Entity, entity, i)
	}
}
