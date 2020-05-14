package typpostgres

const (
	defaultDockerImage = "postgres"
	defaultDockerName  = "postgres"

	defaultConfigName = "PG"

	defaultDBName   = "sample"
	defaultUser     = "postgres"
	defaultPassword = "pgpass"
	defaultHost     = "localhost"
	defaultPort     = 5432

	defaultMigrationSrc = "scripts/db/migration"
	defaultSeedSrc      = "scripts/db/seed"
)

// Settings for postgres
type Settings struct {
	Ctor string

	DockerName  string
	DockerImage string

	ConfigName string

	DBName   string
	User     string
	Password string
	Host     string
	Port     int

	MigrationSrc string
	SeedSrc      string
}

// GetDockerName from setting
func GetDockerName(s *Settings) string {
	if s.DockerName == "" {
		return defaultDockerName
	}
	return s.DockerName
}

// GetDockerImage from setting
func GetDockerImage(s *Settings) string {
	if s.DockerImage == "" {
		return defaultDockerImage
	}
	return s.DockerImage
}

// GetConfigName from setting
func GetConfigName(s *Settings) string {
	if s.ConfigName == "" {
		return defaultConfigName
	}
	return s.ConfigName
}

// GetDBName from setting
func GetDBName(s *Settings) string {
	if s.DBName == "" {
		return defaultDBName
	}
	return s.DBName
}

// GetUser from setting
func GetUser(s *Settings) string {
	if s.User == "" {
		return defaultUser
	}
	return s.User
}

// GetPassword from setting
func GetPassword(s *Settings) string {
	if s.Password == "" {
		return defaultPassword
	}
	return s.Password
}

// GetHost from setting
func GetHost(s *Settings) string {
	if s.Host == "" {
		return defaultHost
	}
	return s.Host
}

// GetPort from setting
func GetPort(s *Settings) int {
	if s.Port == 0 {
		return defaultPort
	}
	return s.Port
}

// GetMigrationSrc from setting if available or the default value
func GetMigrationSrc(s *Settings) string {
	if s.MigrationSrc == "" {
		return defaultMigrationSrc
	}
	return s.MigrationSrc
}

// GetSeedSrc from setting if available or the default value
func GetSeedSrc(s *Settings) string {
	if s.SeedSrc == "" {
		return defaultSeedSrc
	}
	return s.SeedSrc
}
