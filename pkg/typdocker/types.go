package typdocker

import "strings"

// ComposeObject represent docker-compose.yml
type ComposeObject struct {
	Version  Version
	Services Services
	Networks Networks
	Volumes  Volumes
}

// Append another compose object
func (c *ComposeObject) Append(comp *ComposeObject) {
	for k, service := range comp.Services {
		c.Services[k] = service
	}
	for k, network := range comp.Networks {
		c.Networks[k] = network
	}
	for k, volume := range comp.Volumes {
		c.Volumes[k] = volume
	}
}

// ------------------------------------------------------------------

// Services descriptor in docker-compose.yml
type Services map[string]interface{}

// ------------------------------------------------------------------

// Networks descriptor in docker-compose.yml
type Networks map[string]interface{}

// ------------------------------------------------------------------

// Volumes descriptor in docker-compose.yml
type Volumes map[string]interface{}

// ------------------------------------------------------------------

// Network in docker-compose.yaml
type Network struct {
	Driver string `yaml:"driver,omitempty"`
}

// ------------------------------------------------------------------

// Service in docker-compose.yaml
type Service struct {
	Image       string            `yaml:"image,omitempty"`
	Command     string            `yaml:"command,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
}

// ------------------------------------------------------------------

// Version of docker compose
type Version string

// IsV3 return true if the version belong to docker-compose v3
func (v *Version) IsV3() bool {
	// TODO: improve the checking
	return strings.HasPrefix(string(*v), "3")
}
