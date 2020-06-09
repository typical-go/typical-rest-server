package typrails_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typrails"
)

func TestCreateEntity(t *testing.T) {
	testcases := []struct {
		column string
		udt    string

		typrails.Field
	}{
		{
			column: "id",
			udt:    "int4",
			Field: typrails.Field{
				Name:      "ID",
				Type:      "int64",
				Udt:       "int4",
				Column:    "id",
				StructTag: "`json:\"id\"`",
			},
		},
		{
			column: "created_at",
			udt:    "timestamp",
			Field: typrails.Field{
				Name:      "CreatedAt",
				Type:      "time.Time",
				Udt:       "timestamp",
				Column:    "created_at",
				StructTag: "`json:\"created_at\"`",
			},
		},
		{
			column: "name",
			udt:    "varchar",
			Field: typrails.Field{
				Name:      "Name",
				Type:      "string",
				Udt:       "varchar",
				Column:    "name",
				StructTag: "`json:\"name\"`",
			},
		},
	}
	for _, tt := range testcases {
		require.EqualValues(t, tt.Field, typrails.CreateField(tt.column, tt.udt))
	}

}
