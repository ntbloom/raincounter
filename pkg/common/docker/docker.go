package docker

import (
	_ "github.com/docker/docker/client"
)

// DockerContainer wraps the docker api to launch containers for testing purposes
// Not to be used as production code
type DockerContainer struct {
	image string
	name  string
	port  int
}

func NewDockerContainer(image, name string, port int) *DockerContainer {
	return &DockerContainer{image, name, port}
}

// Run launches an ephemeral container
func (d *DockerContainer) Run() {
}

// Kill removes the container
func (d *DockerContainer) Kill() {}
