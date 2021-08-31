package receiver_test

import (
	"testing"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/server/receiver"

	"github.com/stretchr/testify/suite"
)

const localhost = "127.0.0.1"

type ReceiverTest struct {
	suite.Suite
	testDatabase string
	receiver     *receiver.Receiver
}

func TestReceiver(t *testing.T) {
	test := new(ReceiverTest)
	suite.Run(t, test)
}

func (suite *ReceiverTest) SetupSuite() {
	config.Configure()

	// connect to the docker container without auth
	client, err := mqtt.NewConnection(mqtt.NewBrokerConfigNoAuth(localhost, 1883))
	if err != nil {
		panic(err)
	}

	testFile := "deadbeef"
	suite.testDatabase = testFile

	r, err := receiver.NewReceiver(client, testFile)
	if err != nil {
		panic(err)
	}
	suite.receiver = r
}

func (suite *ReceiverTest) TearDownSuite() {}

func (suite *ReceiverTest) SetupTest()    {}
func (suite *ReceiverTest) TearDownTest() {}

// can we actually connect to the mqtt container?
func (suite *ReceiverTest) TestBasicConnection() {
	assert.True(suite.T(), suite.receiver.IsConnected())
}
