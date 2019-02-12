package config

// FIXME: Change application name, usage, version

var (
	// App contain application detail
	App = struct {
		Name    string
		Usage   string
		Version string
	}{
		Name:    "[Name]",
		Usage:   "API for [Usage]",
		Version: "0.0.1",
	}

	// Prefix of configuration
	Prefix = "APP"

	// DefaultMigrationDirectory refer to migration directory path
	DefaultMigrationDirectory = "db/migrate"
)
