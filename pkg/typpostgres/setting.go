package typpostgres

const (
	DefaultDockerImage = "postgres"
	DefaultDockerName  = "postgres"

	DefaultConfigName = "PG"

	DefaultDBName   = "sample"
	DefaultUser     = "postgres"
	DefaultPassword = "pgpass"
	DefaultHost     = "localhost"
	DefaultPort     = 5432

	DefaultMigrationSrc = "scripts/db/migration"
	DefaultSeedSrc      = "scripts/db/seed"
)

// Setting for postgres
type Setting struct {
	CtorName string

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
func GetDockerName(s *Setting) string {
	if s.DockerName == "" {
		return DefaultDockerName
	}
	return s.DockerName
}

// GetDockerImage from setting
func GetDockerImage(s *Setting) string {
	if s.DockerImage == "" {
		return DefaultDockerImage
	}
	return s.DockerImage
}

// GetConfigName from setting
func GetConfigName(s *Setting) string {
	if s.ConfigName == "" {
		return DefaultConfigName
	}
	return s.ConfigName
}

// GetDBName from setting
func GetDBName(s *Setting) string {
	if s.DBName == "" {
		return DefaultDBName
	}
	return s.DBName
}

// GetUser from setting
func GetUser(s *Setting) string {
	if s.User == "" {
		return DefaultUser
	}
	return s.User
}

// GetPassword from setting
func GetPassword(s *Setting) string {
	if s.Password == "" {
		return DefaultPassword
	}
	return s.Password
}

// GetHost from setting
func GetHost(s *Setting) string {
	if s.Host == "" {
		return DefaultHost
	}
	return s.Host
}

// GetPort from setting
func GetPort(s *Setting) int {
	if s.Port == 0 {
		return DefaultPort
	}
	return s.Port
}

// GetMigrationSrc from setting if available or the default value
func GetMigrationSrc(s *Setting) string {
	if s.MigrationSrc == "" {
		return DefaultMigrationSrc
	}
	return s.MigrationSrc
}

// GetSeedSrc from setting if available or the default value
func GetSeedSrc(s *Setting) string {
	if s.SeedSrc == "" {
		return DefaultSeedSrc
	}
	return s.SeedSrc
}
