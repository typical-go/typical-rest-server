package typredis

var (
	// DefaultConfigName is default value for redis config name
	DefaultConfigName = "REDIS"

	// DefaultHost is default value for redis host
	DefaultHost = "localhost"

	// DefaultPort is default value for redis port
	DefaultPort = "6379"

	// DefaultPassword is default value for redis password
	DefaultPassword = "redispass"

	// DefaultConfig for redis
	DefaultConfig = &Config{
		Host:     DefaultHost,
		Port:     DefaultPort,
		Password: DefaultPassword,
	}
)
