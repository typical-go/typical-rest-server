package typdocker

// Network in docker-compose.yaml
type Network struct {
	Driver string `yaml:"driver,omitempty"`
}
