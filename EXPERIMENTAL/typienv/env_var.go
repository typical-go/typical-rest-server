package typienv

import "os"

// EnvVar environment variable
type EnvVar struct {
	Name    string
	Default string
}

// Value of environment variable
func (e EnvVar) Value() (value string) {
	value = os.Getenv(e.Name)
	if value == "" {
		value = e.Default
	}
	return value
}
