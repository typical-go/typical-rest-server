package typiparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAutoWire(t *testing.T) {
	testcase := []struct {
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

	for _, tt := range testcase {
		require.Equal(t, tt.autowire, isAutoWire(tt.funcName, tt.doc))
	}
}

func TestIsAutoMock(t *testing.T) {
	testcase := []struct {
		doc      string
		automock bool
	}{
		{"some doc", true},
		{"some doc [nomock]", true},
	}

	for _, tt := range testcase {
		require.Equal(t, tt.automock, isAutoMock(tt.doc))
	}
}
