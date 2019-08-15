package envkit

import (
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	env := map[string]string{
		"Foo":   "Bar",
		"Hello": "World",
	}
	Set(env)
	defer os.Clearenv()

	for key, value := range env {
		s := os.Getenv(key)
		if value != s {
			t.Fatalf("'%s': want '%s' but got '%s'", key, value, s)
		}
	}
}
