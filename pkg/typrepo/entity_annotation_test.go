package typrepo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/typrepo"
)

func TestGetDest(t *testing.T) {
	testcases := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{Path: "internal/app/data_access/mysql/a.go", Expected: "internal/generated/entity/app/data_access/mysql_repo"},
		{Path: "internal/app/data_access/mysql_entity/a.go", Expected: "internal/generated/entity/app/data_access/mysql_repo"},
		{Path: "internal/app/entity/a.go", Expected: "internal/generated/entity/app/repo"},
	}
	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			a := typrepo.EntityAnnotation{}
			require.Equal(t, tt.Expected, a.GetDest(tt.Path))
		})
	}
}
