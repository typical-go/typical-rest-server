package typdocker

type (
	// Recipe represent docker-compose.yml
	Recipe struct {
		Version  string
		Services Services
		Networks Networks
		Volumes  Volumes
	}
	// Services descriptor in docker-compose.yml
	Services map[string]interface{}
	// Networks descriptor in docker-compose.yml
	Networks map[string]interface{}
	// Volumes descriptor in docker-compose.yml
	Volumes map[string]interface{}
	// Network in docker-compose.yaml
	Network struct {
		Driver string `yaml:"driver,omitempty"`
	}
	// Service in docker-compose.yaml
	Service struct {
		Image       string            `yaml:"image,omitempty"`
		Command     string            `yaml:"command,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Networks    []string          `yaml:"networks,omitempty"`
		Restart     string            `yaml:"restart,omitempty"`
	}
)

//
// Recipe
//

var _ Composer = (*Recipe)(nil)

// ComposeV3 the recipe
func (c *Recipe) ComposeV3() (*Recipe, error) {
	return c, nil
}
