package config

// AppConfig contain applicatoin configuration
type AppConfig struct {
	Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
}
