package envkit

import "os"

// Set multiple value to environment variable
func Set(env map[string]string) (err error) {
	for key, value := range env {
		err = os.Setenv(key, value)
	}
	return
}
