package typredis

const (
	defaultCmd         = "redis"
	defaultConfigName  = "REDIS"
	defaultHost        = "localhost"
	defaultPort        = "6379"
	defaultPassword    = "redispass"
	defaultDockerName  = "redis"
	defaultDockerImage = "redis:4.0.5-alpine"
)

// Settings of redis
type Settings struct {
	Ctor       string
	UtilityCmd string

	ConfigName string
	Host       string
	Port       string
	Password   string

	DockerName  string
	DockerImage string
}

// Init settings
func Init(s *Settings) *Settings {
	if s.UtilityCmd == "" {
		s.UtilityCmd = defaultCmd
	}
	if s.ConfigName == "" {
		s.ConfigName = defaultConfigName
	}
	if s.Host == "" {
		s.Host = defaultHost
	}
	if s.Port == "" {
		s.Port = defaultPort
	}
	if s.Password == "" {
		s.Password = defaultPassword
	}
	if s.DockerName == "" {
		s.DockerName = defaultDockerName
	}
	if s.DockerImage == "" {
		s.DockerImage = defaultDockerImage
	}
	return s
}
