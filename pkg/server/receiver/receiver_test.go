package receiver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/server/receiver"

	"github.com/stretchr/testify/suite"
)

type ReceiverTest struct {
	suite.Suite
	testFile string
	receiver *receiver.Receiver
}

func TestReceiver(t *testing.T) {
	test := new(ReceiverTest)
	suite.Run(t, test)
}

func (suite *ReceiverTest) SetupSuite() {
	config.Configure()
	suite.testFile = viper.GetString(configkey.DatabaseRemoteDevFile)
	r, err := receiver.NewReceiver(suite.testFile, true)
	if err != nil {
		panic(err)
	}
	suite.receiver = r
}

func (suite *ReceiverTest) SetupTest()     {}
func (suite *ReceiverTest) TearDownTest()  {}
func (suite *ReceiverTest) TearDownSuite() {}

func (suite *ReceiverTest) TestBasicConnection() {
	assert.True(suite.T(), suite.receiver.IsConnected())
}
