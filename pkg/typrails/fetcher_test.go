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
		data map[string]string
		err  error
	}{
		{
			data: map[string]string{
				"column1": "type1",
				"column2": "type2",
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing; \"updated_at\" with underlying data type \"timestamp\" is missing; \"created_at\" with underlying data type \"timestamp\" is missing"),
		},
		{
			data: map[string]string{
				"id":         "int4",
				"created_at": "timestamp",
				"updated_at": "timestamp",
			},
		},
		{
			data: map[string]string{
				"id":         "int",
				"created_at": "timestamp",
				"updated_at": "timestamp",
			},
			err: errors.New("\"id\" with underlying data type \"int4\" is missing"),
		},
	}
	fetcher := typrails.Fetcher{DB: db}
	query := regexp.QuoteMeta("SELECT column_name, udt_name FROM information_schema.COLUMNS WHERE table_name = ?")
	for _, tt := range testcases {
		rows := sqlmock.NewRows([]string{"column_name", "data_type"})
		for key, value := range tt.data {
			rows.AddRow(key, value)
		}
		mock.ExpectQuery(query).WithArgs("some-table").WillReturnRows(rows)
		_, err = fetcher.Fetch("some-table")
		require.EqualValues(t, tt.err, err)
	}
}
