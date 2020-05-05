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

// Setting for postgres
type Setting struct {
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

// DockerName from setting
func DockerName(s *Setting) string {
	if s.DockerName == "" {
		return defaultDockerName
	}
	return s.DockerName
}

// DockerImage from setting
func DockerImage(s *Setting) string {
	if s.DockerImage == "" {
		return defaultDockerImage
	}
	return s.DockerImage
}

// ConfigName from setting
func ConfigName(s *Setting) string {
	if s.ConfigName == "" {
		return defaultConfigName
	}
	return s.ConfigName
}

// DBName from setting
func DBName(s *Setting) string {
	if s.DBName == "" {
		return defaultDBName
	}
	return s.DBName
}

// User from setting
func User(s *Setting) string {
	if s.User == "" {
		return defaultUser
	}
	return s.User
}

// Password from setting
func Password(s *Setting) string {
	if s.Password == "" {
		return defaultPassword
	}
	return s.Password
}

// Host from setting
func Host(s *Setting) string {
	if s.Host == "" {
		return defaultHost
	}
	return s.Host
}

// Port from setting
func Port(s *Setting) int {
	if s.Port == 0 {
		return defaultPort
	}
	return s.Port
}

// MigrationSrc from setting
func MigrationSrc(s *Setting) string {
	if s.MigrationSrc == "" {
		return defaultMigrationSrc
	}
	return s.MigrationSrc
}

// SeedSrc from setting
func SeedSrc(s *Setting) string {
	if s.SeedSrc == "" {
		return defaultSeedSrc
	}
	return s.SeedSrc
}
