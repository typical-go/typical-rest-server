package typdocker

import "strings"

// Version of docker compose
type Version string

// IsV3 return true if the version belong to docker-compose v3
func (v *Version) IsV3() bool {
	// TODO: improve the checking
	return strings.HasPrefix(string(*v), "3")
}
