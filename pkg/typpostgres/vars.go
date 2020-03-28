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

	// DefaultDockerImage is default docker image for postgres
	DefaultDockerImage = "postgres"

	// DefaultDockerName is default docker name for postgres
	DefaultDockerName = "postgres"

	// DefaultMigrationSource is default migration source for postgres
	DefaultMigrationSource = "scripts/db/migration"

	// DefaultSeedSource is default seed source for postgres
	DefaultSeedSource = "scripts/db/seed"
)
