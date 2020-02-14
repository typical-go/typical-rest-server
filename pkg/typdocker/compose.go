package typdocker

func (m *Module) dockerCompose() (root *ComposeObject) {
	root = &ComposeObject{
		Version:  m.Version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, composer := range m.Composers {
		if obj := composer.DockerCompose(m.Version); obj != nil {
			root.Append(obj)
		}
	}
	return
}
