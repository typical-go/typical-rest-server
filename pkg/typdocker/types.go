package typdocker

// DockerComposer responsible to compose docker
type DockerComposer interface {
	DockerCompose() Compose
}

// Network in docker-compose.yaml
type Network struct {
	Driver string `yaml:"driver,omitempty"`
}

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

// Compose is abstraction for docker-compose.yml
type Compose struct {
	Version  string
	Services map[string]interface{}
	Networks map[string]interface{}
	Volumes  map[string]interface{}
}

// Add another compose
func (c *Compose) Add(comp Compose) {
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
