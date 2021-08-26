package docker

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	_ "github.com/docker/docker/client"
)

// DockerContainer wraps the docker api to launch containers for testing purposes only.
// Uses the docker SDK to run the equivalent of `docker run --rm -d -p 8080:8080 --name my_container nginx`.
type DockerContainer struct {
	image         string
	name          string
	hostPort      int
	containerPort int
	ctx           context.Context
	client        *client.Client
	id            string
}

func NewDockerContainer(image, name string, hostPort, containerPort int) (*DockerContainer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &DockerContainer{image, name, hostPort, containerPort, context.Background(), cli, ""}, nil
}

// Run launches an ephemeral container
func (d *DockerContainer) Run() error {
	if err := d.pull(); err != nil {
		return err
	}
	if err := d.create(); err != nil {
		return err
	}
	if err := d.start(); err != nil {
		_ = d.forceRemove()
		return err
	}
	return nil
}

// Kill removes the container
func (d *DockerContainer) Kill() error {
	return d.forceRemove()
}

// pull the latest image
func (d *DockerContainer) pull() error {
	out, err := d.client.ImagePull(d.ctx, d.image, types.ImagePullOptions{})
	if err != nil {
		logrus.Error(err)
		return err
	}
	if _, err := io.Copy(os.Stderr, out); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// create the container
func (d *DockerContainer) create() error {
	port, err := nat.NewPort("tcp", strconv.Itoa(d.containerPort))
	if err != nil {
		logrus.Error(nil)
		return err
	}
	containerCfg := &container.Config{
		Image: d.image,
		ExposedPorts: nat.PortSet{
			port: struct{}{},
		},
	}
	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(d.hostPort),
				},
			},
		},
	}
	resp, err := d.client.ContainerCreate(
		d.ctx,
		containerCfg,
		hostCfg,
		nil,
		nil,
		d.name,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	d.id = resp.ID
	return nil
}

// start the container
func (d *DockerContainer) start() error {
	if d.id == "" {
		panic("container ID not set, did you pull and create the container first?")
	}

	options := types.ContainerStartOptions{}
	if err := d.client.ContainerStart(d.ctx, d.id, options); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// forceRemove the container
func (d *DockerContainer) forceRemove() error {
	if d.id == "" {
		panic("container ID not set, did you pull and create the container first?")
	}
	options := types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	}
	if err := d.client.ContainerRemove(d.ctx, d.id, options); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}