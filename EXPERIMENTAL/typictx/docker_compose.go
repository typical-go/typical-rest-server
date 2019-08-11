package typictx

// DockerCompose is abstraction for docker-compose.yml
type DockerCompose struct {
	Version  string
	Services map[string]interface{}
	Networks map[string]interface{}
	Volumes  map[string]interface{}
}
