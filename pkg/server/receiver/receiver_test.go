package receiver_test

import (
	"testing"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/common/docker"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/server/receiver"

	"github.com/stretchr/testify/suite"
)

const localhost = "127.0.0.1"

type ReceiverTest struct {
	suite.Suite
	testFile  string
	receiver  *receiver.Receiver
	mosquitto *docker.Container
}

func TestReceiver(t *testing.T) {
	test := new(ReceiverTest)
	suite.Run(t, test)
}

func (suite *ReceiverTest) SetupSuite() {
	config.Configure()

	// launch the docker container
	container, err := docker.NewContainer("eclipse-mosquitto", "receiver-test", 1883)
	if err != nil {
		panic(err)
	}
	suite.mosquitto = container
	if err = suite.mosquitto.Run(); err != nil {
		panic(err)
	}

	// connect to the docker container without auth
	client, err := mqtt.NewConnection(mqtt.NewBrokerConfigNoAuth(localhost, 1883))
	if err != nil {
		panic(err)
	}

	// prep the localdb file
	testFile := viper.GetString(configkey.DatabaseRemoteFile)
	suite.testFile = testFile

	r, err := receiver.NewReceiver(client, testFile, true)
	if err != nil {
		panic(err)
	}
	suite.receiver = r
}

func (suite *ReceiverTest) TearDownSuite() {
	suite.mosquitto.Kill()
}

func (suite *ReceiverTest) SetupTest()    {}
func (suite *ReceiverTest) TearDownTest() {}

// can we actually connect to the mqtt container?
func (suite *ReceiverTest) TestBasicConnection() {
	assert.True(suite.T(), suite.receiver.IsConnected())
}
