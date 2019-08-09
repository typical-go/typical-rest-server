package config

// Config is application configuration
type Config struct {
	Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
}
