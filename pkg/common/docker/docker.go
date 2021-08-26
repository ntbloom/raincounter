package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	_ "github.com/docker/docker/client"
)

// DockerContainer wraps the docker api to launch containers for testing purposes only
type DockerContainer struct {
	image         string
	name          string
	hostPort      int
	containerPort int
	ctx           context.Context
	cli           *client.Client
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

// Pull fetches the latest image, maybe not necessary?
func (d *DockerContainer) Pull() error {
	out, err := d.cli.ImagePull(d.ctx, d.image, types.ImagePullOptions{})
	if err != nil {
		logrus.Error(err)
		return err
	}
	if _, err := io.Copy(os.Stdout, out); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// Run launches an ephemeral container
func (d *DockerContainer) Run() error {
	if err := d.Pull(); err != nil {
		return err
	}
	containerCfg := &container.Config{
		Image: d.image,
		ExposedPorts: nat.PortSet{
			"5432:tcp": struct{}{},
		},
	}
	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "5432",
				},
			},
		},
	}
	resp, err := d.cli.ContainerCreate(
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

// Kill removes the container
func (d *DockerContainer) Kill() error {
	return nil
}
