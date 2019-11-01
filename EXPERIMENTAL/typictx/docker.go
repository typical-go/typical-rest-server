package typictx

// DockerCompose get docker compose
// func (c *Context) DockerCompose() (dockerCompose *docker.Compose) {
// 	dockerCompose = docker.NewCompose("3")
// 	for _, module := range c.Modules {
// 		moduleDocker := module.DockerCompose
// 		if moduleDocker == nil {
// 			continue
// 		}
// 		for _, name := range moduleDocker.ServiceKeys {
// 			dockerCompose.RegisterService(name, moduleDocker.Services[name])
// 		}
// 		for _, name := range moduleDocker.NetworkKeys {
// 			dockerCompose.RegisterNetwork(name, moduleDocker.Networks[name])
// 		}
// 		for _, name := range moduleDocker.VolumeKeys {
// 			dockerCompose.RegisterVolume(name, moduleDocker.Volumes[name])
// 		}
// 	}
// 	return
// }
