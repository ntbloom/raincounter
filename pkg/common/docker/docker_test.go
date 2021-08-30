package docker_test

import (
	"os"
	"sort"
	"time"

	"github.com/ntbloom/raincounter/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/ntbloom/raincounter/pkg/common/docker"
)

const (
	localhost = "127.0.0.1"
	image     = "postgres"
)

type DockerTest struct {
	suite.Suite
}

func TestPostgresql(t *testing.T) {
	test := new(DockerTest)
	suite.Run(t, test)
}

func (suite *DockerTest) SetupSuite() {
	config.Configure()
}
func (suite *DockerTest) TearDownSuite() {}
func (suite *DockerTest) SetupTest()     {}
func (suite *DockerTest) TearDownTest()  {}

// can we actually connect to the mqtt container?
func (suite *DockerTest) TestPostgresql() {
	container, err := docker.NewContainer(image, "postgresql-test", 5432)
	if err != nil {
		suite.Fail("container struct not instantiated")
	}
	defer container.Kill()
	if err = os.Setenv("POSTGRES_PASSWORD", "password"); err != nil {
		logrus.Error("problem setting ENV variables")
		suite.Fail("not propagating environment variables")
	}
	env := os.Environ()
	sort.Strings(env)
	logrus.Trace(env)
	if err = container.Run(); err != nil {
		logrus.Error(err)
		suite.Fail("container not running")
	}
	time.Sleep(time.Second * 5)
	assert.True(suite.T(), container.Status() == "running", "container is not running")

}
