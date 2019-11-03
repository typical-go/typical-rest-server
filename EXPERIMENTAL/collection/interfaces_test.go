package collection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/collection"
)

func TestInterfaces(t *testing.T) {
	var coll collection.Interfaces
	coll.Add("some-item")
	coll.Add(88)
	coll.Add(3.14)
	require.EqualValues(t, []interface{}{"some-item", 88, 3.14}, coll)
}
