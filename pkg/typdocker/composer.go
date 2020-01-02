package typdocker

// Composer responsible to compose docker
type Composer interface {
	DockerCompose(version Version) *ComposeObject
}

// ComposeObject represent docker-compose.yml
type ComposeObject struct {
	Version  Version
	Services Services
	Networks Networks
	Volumes  Volumes
}

// Services descriptor in docker-compose.yml
type Services map[string]interface{}

// Networks descriptor in docker-compose.yml
type Networks map[string]interface{}

// Volumes descriptor in docker-compose.yml
type Volumes map[string]interface{}

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
