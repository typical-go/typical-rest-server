package typrails_test

import (
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
	fetcher := typrails.Fetcher{
		DB: db,
	}
	query := regexp.QuoteMeta("SELECT column_name, data_type FROM information_schema.COLUMNS WHERE table_name = ?")
	t.Run("WHEN invalid table", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs("some-table").
			WillReturnRows(sqlmock.NewRows([]string{"column_name", "data_type"}).
				AddRow("column1", "type1").
				AddRow("column2", "type2"))
		_, err = fetcher.Fetch("some-table")
		require.EqualError(t, err, "\"id\" is missing; \"updated_at\" is missing; \"created_at\" is missing")
	})

}
