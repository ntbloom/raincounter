package receiver_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"
	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/server/receiver"
	"github.com/ntbloom/raincounter/pkg/server/webdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const localhost = "127.0.0.1"

type ReceiverTest struct {
	suite.Suite
	receiver *receiver.Receiver
	client   paho.Client
	query    webdb.DBQuery
	entry    webdb.DBEntry
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
		suite.Fail("unable to connect to mqtt", err)
	}
	suite.client = client
	r, err := receiver.NewReceiver(client)
	if err != nil {
		suite.Fail("unable to make a new Receiver struct", err)
	}
	suite.receiver = r

	// control the database as well
	var query webdb.DBQuery
	var entry webdb.DBEntry
	db := webdb.NewPGConnector()
	query = db
	entry = db
	suite.query = query
	suite.entry = entry
}

func (suite *ReceiverTest) TearDownSuite() {
	logrus.Debug("closing the database pool")
	suite.query.Close()
	suite.entry.Close()
	logrus.Debug("disconnecting the client from mqtt")
	suite.client.Disconnect(viper.GetUint(configkey.MQTTQuiescence))
	logrus.Debug("disconnecting test receiver from mqtt")
	suite.receiver.Close()
}

func (suite *ReceiverTest) SetupTest() {
	logrus.Info("deleting all rows from receiver_test.go")
	// delete all database rows
	for _, sql := range []string{
		"DELETE FROM temperature;",
		"DELETE FROM rain;",
		"DELETE FROM event_log;",
		"DELETE FROM status_log;",
	} {
		// `Select` can still execute arbitrary SQL
		err := suite.entry.Insert(sql)
		if err != nil {
			logrus.Error(err)
			suite.Fail("can't delete table rows", err)
		}
	}
}
func (suite *ReceiverTest) TearDownTest() {}

// publish a rain topic, make sure it gets into the database
func (suite *ReceiverTest) TestReceiveRainMessage() {
	msg := mqtt.SampleRain()
	suite.client.Publish(process(msg))
	// wait for it to make it to the broker
	time.Sleep(time.Second * 1)

	// verify the last rain matches what we put in the database
	lastRain, err := suite.query.GetLastRainTime()
	if err != nil {
		suite.Fail("last rain error", err)
	}
	timeDiff := msg.Timestamp.Sub(lastRain)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	assert.True(suite.T(), timeDiff < time.Second)
}

func (suite *ReceiverTest) TestReceiveTemperatureMessage() {
	msg := mqtt.SampleTemp()
	suite.client.Publish(process(msg))
	lastTemp, err := suite.query.GetLastTempC()
	if err != nil {
		suite.Fail("last temperature error", err)
	}
	expTemp := msg.Msg["TempC"]
	assert.Equal(suite.T(), lastTemp, expTemp)
}

// TODO: come back to this test
//func (suite *ReceiverTest) TestAllMessages() {
//	for _, message := range []mqtt.SampleMessage{
//		mqtt.SampleSensorPause,
//		mqtt.SampleSensorUnpause,
//		mqtt.SampleSensorSoftReset,
//		mqtt.SampleSensorHardReset,
//		mqtt.SampleSensorStatus,
//		mqtt.SampleGatewayStatus,
//	} {
//		suite.client.Publish(process(message))
//	}
//
//	// wait for it to make it to the broker
//	time.Sleep(time.Second)
//
//	// verify the last rain matches what we put in the database
//	panic("implement me!")
//}

// publish a bunch of stuff to the broker
func process(msg mqtt.SampleMessage) (string, byte, bool, []byte) {
	payload, err := json.Marshal(msg.Msg)
	if err != nil {
		logrus.Error(err)
		panic("problem marshalling json")
	}
	return msg.Topic, 1, false, payload
}
