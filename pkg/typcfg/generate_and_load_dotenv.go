package typcfg

import (
	"strings"

	"github.com/typical-go/typical-go/pkg/envkit"
)

// GenerateAndLoadDotEnv to create and load envfile
func GenerateAndLoadDotEnv(target string, c *Context) error {
	envmap, err := envkit.ReadFile(target)
	if err != nil {
		envmap = make(map[string]string)
	}

	var updatedKeys []string
	for _, Envconfig := range c.Configs {
		for _, field := range Envconfig.Fields {
			if _, ok := envmap[field.Key]; !ok {
				updatedKeys = append(updatedKeys, field.Key)
				envmap[field.Key] = field.Default
			}
		}
	}
	if len(updatedKeys) > 0 {
		c.Infof("New keys added in '%s': %s\n", target, strings.Join(updatedKeys, " "))
	}

	if err := envkit.SaveFile(envmap, target); err != nil {
		return err
	}

	return envkit.Setenv(envmap)
}
