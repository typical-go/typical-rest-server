package typredis

import "os"

var DefaultEnvKeys = &EnvKeys{
	Host: "HOST",
	Port: "PORT",
	Pass: "PASS",
}

type (
	Config struct {
		Host string
		Port string
		Pass string
	}
	EnvKeys Config
)

func EnvKeysWithPrefix(prefix string) *EnvKeys {
	return &EnvKeys{
		Host: prefix + "_" + DefaultEnvKeys.Host,
		Port: prefix + "_" + DefaultEnvKeys.Port,
		Pass: prefix + "_" + DefaultEnvKeys.Pass,
	}
}

func (r *EnvKeys) Config() *Config {
	return &Config{
		Host: os.Getenv(r.Host),
		Port: os.Getenv(r.Port),
		Pass: os.Getenv(r.Pass),
	}
}
