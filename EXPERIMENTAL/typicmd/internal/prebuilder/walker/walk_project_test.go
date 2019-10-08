package walker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAutoWire(t *testing.T) {
	testcases := []struct {
		funcName string
		doc      string
		autowire bool
	}{
		{"NewService", "some doc", true},
		{"SomeFunction", "some doc", false},
		{"SomeFunction", "some doc [autowire]", true},
		{"SomeFunction", "some doc [Autowire]", true},
		{"SomeFunction", "some doc [AUTOWIRE]", true},
		{"NewSomething", "some doc [nowire]", false},
		{"NewSomething", "some doc [nowire][autowrite]", false},
		{"Something", "some doc [nowire][autowire]", true},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.autowire, isAutoWire(tt.funcName, tt.doc))
	}
}

func TestIsAutoMock(t *testing.T) {
	testcases := []struct {
		doc      string
		automock bool
	}{
		{"some doc", true},
		{"some doc [nomock]", false},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.automock, isAutoMock(tt.doc))
	}
}

func TestWalkTarget(t *testing.T) {
	testcases := []struct {
		filename string
		result   bool
	}{
		{"file.go", true},
		{"file_test.go", false},
		{"file.test.go", true},
		{"file", false},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.result, isWalkTarget(tt.filename))
	}
}
