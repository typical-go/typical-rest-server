package dbtool

type (
	// Config config for postgres tool
	Config struct {
		DBName string
		DBUser string
		DBPass string
		Host   string
		Port   string
	}
	// Configurer return config
	Configurer interface {
		Config() *Config
	}
)
