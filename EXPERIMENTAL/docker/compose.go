package docker

// Compose is abstraction for docker-compose.yml
type Compose struct {
	Version     string
	Services    map[string]interface{}
	Networks    map[string]interface{}
	Volumes     map[string]interface{}
	ServiceKeys []string `yaml:"-"`
	NetworkKeys []string `yaml:"-"`
	VolumeKeys  []string `yaml:"-"`
}

// NewCompose return new instance of DockerCompose
func NewCompose(version string) *Compose {
	return &Compose{
		Version: version,
	}
}

// RegisterService to register service
func (d *Compose) RegisterService(name string, service interface{}) *Compose {
	if d.Services == nil {
		d.Services = make(map[string]interface{})
	}
	_, ok := d.Services[name]
	if !ok {
		d.ServiceKeys = append(d.ServiceKeys, name)
	}
	d.Services[name] = service
	return d
}

// RegisterNetwork to register network
func (d *Compose) RegisterNetwork(name string, network interface{}) *Compose {
	if d.Networks == nil {
		d.Networks = make(map[string]interface{})
	}
	_, ok := d.Networks[name]
	if !ok {
		d.NetworkKeys = append(d.NetworkKeys, name)
	}
	d.Networks[name] = network
	return d
}

// RegisterVolume to register volume
func (d *Compose) RegisterVolume(name string, volume interface{}) *Compose {
	if d.Volumes == nil {
		d.Volumes = make(map[string]interface{})
	}
	_, ok := d.Volumes[name]
	if !ok {
		d.VolumeKeys = append(d.VolumeKeys, name)
	}
	d.Volumes[name] = volume
	return d
}
