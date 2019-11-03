package docker

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
