package receiver_test

import (
	"testing"

	"github.com/ntbloom/raincounter/pkg/common/docker"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/server/receiver"

	"github.com/stretchr/testify/suite"
)

type ReceiverTest struct {
	suite.Suite
	testFile  string
	receiver  *receiver.Receiver
	mosquitto *docker.DockerContainer
}

func TestReceiver(t *testing.T) {
	test := new(ReceiverTest)
	suite.Run(t, test)
}

func (suite *ReceiverTest) SetupSuite() {
	config.Configure()

	// launch the docker container
	container, err := docker.NewDockerContainer("eclipse-mosquitto", "receiver-test", 1883)
	if err != nil {
		panic(err)
	}
	suite.mosquitto = container
	if err = suite.mosquitto.Run(); err != nil {
		panic(err)
	}

	// prep the Receiver struct
	suite.testFile = viper.GetString(configkey.DatabaseRemoteDevFile)
	mqttConfig := mqtt.NewBrokerConfig()
	mqttConfig.SetDevBroker("127.0.0.1", 1883)
	r, err := receiver.NewReceiver(mqttConfig, suite.testFile, true)
	if err != nil {
		panic(err)
	}
	suite.receiver = r
}

func (suite *ReceiverTest) SetupTest() {}
func (suite *ReceiverTest) TearDownTest() {

}
func (suite *ReceiverTest) TearDownSuite() {
	if err := suite.mosquitto.Kill(); err != nil {
		panic(err)
	}
}

func (suite *ReceiverTest) TestBasicConnection() {
	assert.True(suite.T(), suite.receiver.IsConnected())
}
