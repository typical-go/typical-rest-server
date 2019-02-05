package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("TEST_ADDRESS", ":8889")
	defer os.Clearenv()

	conf, err := Load("TEST")

	assert.NoError(t, err)
	assert.Equal(t, conf.Address, ":8889")
}
