package typictx

// DockerCompose is abstraction for docker-compose.yml
type DockerCompose struct {
	Version  string
	Services map[string]interface{}
	Networks map[string]interface{}
	Volumes  map[string]interface{}
}

// NewDockerCompose return new instance of DockerCompose
func NewDockerCompose(version string) *DockerCompose {
	return &DockerCompose{
		Version: version,
	}
}

// RegisterService to register service
func (d *DockerCompose) RegisterService(name string, service interface{}) *DockerCompose {
	if d.Services == nil {
		d.Services = make(map[string]interface{})
	}
	d.Services[name] = service
	return d
}

// RegisterNetwork to register network
func (d *DockerCompose) RegisterNetwork(name string, network interface{}) *DockerCompose {
	if d.Networks == nil {
		d.Networks = make(map[string]interface{})
	}
	d.Networks[name] = network
	return d
}

// RegisterVolume to register volume
func (d *DockerCompose) RegisterVolume(name string, volume interface{}) *DockerCompose {
	if d.Volumes == nil {
		d.Volumes = make(map[string]interface{})
	}
	d.Volumes[name] = volume
	return d
}
