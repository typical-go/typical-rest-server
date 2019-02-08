package config_test

import (
	"os"
	"testing"

	"github.com/imantung/typical-go-server/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("TEST_ADDRESS", ":8889")
	defer os.Clearenv()

	config.Prefix = "TEST"
	conf, err := config.NewConfig()

	assert.NoError(t, err)
	assert.Equal(t, conf.Address, ":8889")
}
