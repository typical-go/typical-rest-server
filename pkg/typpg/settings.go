package typpg

const (
	defaultDockerImage = "postgres"
	defaultDockerName  = "postgres"

	defaultConfigName = "PG"
	defaultUtilityCmd = "pg"

	defaultDBName   = "sample"
	defaultUser     = "postgres"
	defaultPassword = "pgpass"
	defaultHost     = "localhost"
	defaultPort     = 5432

	defaultMigrationSrc = "scripts/db/migration"
	defaultSeedSrc      = "scripts/db/seed"
)

type (
	// Settings for postgres
	Settings struct {
		Ctor       string
		UtilityCmd string

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
)

// Init settings
func Init(s *Settings) *Settings {
	if s.UtilityCmd == "" {
		s.UtilityCmd = defaultUtilityCmd
	}
	if s.DockerName == "" {
		s.DockerName = defaultDockerName
	}
	if s.DockerImage == "" {
		s.DockerImage = defaultDockerImage
	}
	if s.ConfigName == "" {
		s.ConfigName = defaultConfigName
	}
	if s.DBName == "" {
		s.DBName = defaultDBName
	}
	if s.User == "" {
		s.User = defaultUser
	}
	if s.Password == "" {
		s.Password = defaultPassword
	}
	if s.Host == "" {
		s.Host = defaultHost
	}
	if s.Port == 0 {
		s.Port = defaultPort
	}
	if s.MigrationSrc == "" {
		s.MigrationSrc = defaultMigrationSrc
	}
	if s.SeedSrc == "" {
		s.SeedSrc = defaultSeedSrc
	}
	return s
}
