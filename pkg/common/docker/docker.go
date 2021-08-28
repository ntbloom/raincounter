package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Container wraps the docker api to launch containers for testing purposes only.
type Container struct {
	image  string
	name   string
	port   int
	ctx    context.Context
	client *client.Client
	id     string
}

func NewContainer(image, name string, port int) (*Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &Container{image, name, port, context.Background(), cli, ""}, nil
}

// Run launches an ephemeral container
func (c *Container) Run() error {
	if err := c.pull(); err != nil {
		return err
	}
	if err := c.create(); err != nil {
		return err
	}
	if err := c.start(); err != nil {
		_ = c.forceRemove()
		return err
	}
	return nil
}

// Kill removes the container
func (c *Container) Kill() error {
	return c.forceRemove()
}

// pull the latest image
func (c *Container) pull() error {
	out, err := c.client.ImagePull(c.ctx, c.image, types.ImagePullOptions{})
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
func (c *Container) create() error {
	port, err := nat.NewPort("tcp", strconv.Itoa(c.port))
	if err != nil {
		logrus.Error(nil)
		return err
	}
	containerCfg := &container.Config{
		Image: c.image,
		ExposedPorts: nat.PortSet{
			port: struct{}{},
		},
	}
	hostCfg := &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: "host",
	}
	resp, err := c.client.ContainerCreate(
		c.ctx,
		containerCfg,
		hostCfg,
		nil,
		nil,
		c.name,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	c.id = resp.ID
	return nil
}

// start the container
func (c *Container) start() error {
	if c.id == "" {
		panic("container ID not set, did you pull and create the container first?")
	}

	options := types.ContainerStartOptions{}
	if err := c.client.ContainerStart(c.ctx, c.id, options); err != nil {
		logrus.Error(err)
		return err
	}
	// block until container is up
	isRunning := func() bool {
		up := false
		containers, _ := c.client.ContainerList(c.ctx, types.ContainerListOptions{})
		for _, v := range containers {
			if v.State == "running" {
				up = true
			}
		}
		return up
	}
	for i := 0; i < 10; i++ {
		if isRunning() {
			break
		}
		time.Sleep(time.Millisecond * 300) //nolint:gomnd
	}
	if !isRunning() {
		return fmt.Errorf("container not started")
	}
	return nil
}

// forceRemove the container
func (c *Container) forceRemove() error {
	if c.id == "" {
		panic("container ID not set, did you pull and create the container first?")
	}
	options := types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	}
	if err := c.client.ContainerRemove(c.ctx, c.id, options); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
