package typcfg

import (
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/common"
)

// GenerateAndLoadDotEnv to create and load envfile
func GenerateAndLoadDotEnv(target string, c *Context) error {
	envmap, err := common.CreateEnvMapFromFile(target)
	if err != nil {
		envmap = make(common.EnvMap)
	}

	var updatedKeys []string
	for _, AppCfg := range c.Configs {
		for _, field := range AppCfg.Fields {
			if _, ok := envmap[field.Key]; !ok {
				updatedKeys = append(updatedKeys, field.Key)
				envmap[field.Key] = field.Default
			}
		}
	}
	if len(updatedKeys) > 0 {
		fmt.Fprintf(Stdout, "New keys added in '%s': %s\n", target, strings.Join(updatedKeys, " "))
	}

	if err := envmap.SaveToFile(target); err != nil {
		return err
	}

	return common.Setenv(envmap)
}
