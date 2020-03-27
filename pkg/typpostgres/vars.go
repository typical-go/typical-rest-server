package typpostgres

var (
	// DefaultConfigName is default lookup key for postgres configuration
	DefaultConfigName = "PG"

	// DefaultDBName is default value for dbName
	DefaultDBName = "sample"

	// DefaultUser is default value for user
	DefaultUser = "postgres"

	// DefaultPassword is default value for password
	DefaultPassword = "pgpass"

	// DefaultHost is default value for host
	DefaultHost = "localhost"

	// DefaultPort is default value for port
	DefaultPort = 5432

	DefaultDockerImage = "postgres"
	DefaultDockerName  = "postgres"

	defaultMigrationSource = "scripts/db/migration"
	defaultSeedSource      = "scripts/db/seed"
)
