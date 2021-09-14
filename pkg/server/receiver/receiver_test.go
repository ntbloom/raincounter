package receiver_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

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
	client, err := mqtt.NewConnection()
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
	stamp := time.Now().Add(time.Minute * -1)
	msg := mqtt.SampleRain(stamp)
	suite.client.Publish(process(msg))
	// wait for it to make it to the broker
	time.Sleep(time.Second * 1)

	// verify the last rain matches what we put in the database
	lastRain, err := suite.query.GetLastRainTime()
	if err != nil {
		suite.Fail("last rain error", err)
	}
	timeDiff := stamp.Sub(lastRain)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	logrus.Infof("timeDiff:%s, stamp:%s, lastRain:%s", timeDiff, stamp, lastRain)
	assert.True(suite.T(), timeDiff < time.Minute*2, "time mismatch on rain message")
}

func (suite *ReceiverTest) TestReceiveTemperatureMessage() {
	msg := mqtt.SampleTemp(time.Now().Add(time.Minute * -1))
	suite.client.Publish(process(msg))
	// wait for it
	time.Sleep(time.Second)

	lastTemp, err := suite.query.GetLastTempC()
	if err != nil {
		suite.Fail("last temperature error", err)
	}
	expTemp := msg.Msg["TempC"]
	assert.Equal(suite.T(), expTemp, lastTemp)
}

func (suite *ReceiverTest) TestStatusMessages() {
	duration := time.Minute * 5
	now := time.Now()

	// assert that the sensor and gateway are not up
	gwUp, err := suite.query.IsGatewayUp(duration)
	if err != nil {
		suite.Fail("unhandled error on empty IsGatewayUp", err)
	}
	sensorUp, err := suite.query.IsSensorUp(duration)
	if err != nil {
		suite.Fail("unhandled error on empty IsSensorUp", err)
	}
	assert.False(suite.T(), sensorUp, "sensor should not be up")
	assert.False(suite.T(), gwUp, "gateway should not be up")

	// publish the messages and wait for a second
	suite.client.Publish(process(mqtt.SampleSensorStatus(now)))
	suite.client.Publish(process(mqtt.SampleGatewayStatus(now)))
	time.Sleep(time.Second)

	// verify the items were put into the database
	gwUp, err = suite.query.IsGatewayUp(duration)
	if err != nil {
		suite.Fail("error querying gateway is up", err)
	}
	sensorUp, err = suite.query.IsSensorUp(duration)
	if err != nil {
		suite.Fail("error querying sensor is up", err)
	}
	assert.True(suite.T(), gwUp, "gateway should be reporting as up")
	assert.True(suite.T(), sensorUp, "sensor should be reporting as up")
}

// make sure we can handle a sensor event
func (suite *ReceiverTest) TestSensorEvent() {
	testEvent := tlv.SoftReset
	testValue := tlv.SoftResetValue
	testPayload := mqtt.SampleSensorSoftReset
	testTimestamp := time.Now().Add(time.Minute * -5)

	// verify there aren't any events yet
	longTime := time.Now().Add(time.Hour * 24 * 365 * -100)
	res, err := suite.query.GetEventMessagesSince(testEvent, longTime)
	if err != nil {
		suite.Fail("problem querying empty event messages", err)
	}
	assert.Nil(suite.T(), *res)

	// add an event over mqtt
	suite.client.Publish(process(testPayload(testTimestamp)))
	time.Sleep(time.Second)

	// verify it's in the database
	res, err = suite.query.GetEventMessagesSince(testEvent, longTime)
	if err != nil {
		suite.Fail("problem querying event messages", err)
	}
	assert.Equal(suite.T(), 1, len(*res), "should have received one and only one message")
	actualEntry := (*res)[0]
	timeDiff := actualEntry.Timestamp.Sub(testPayload(testTimestamp).Timestamp)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	assert.True(suite.T(), timeDiff < time.Second, "mismatched timestamps")
	assert.Equal(suite.T(), testEvent, actualEntry.Tag, "mismatched tag")
	assert.Equal(suite.T(), testValue, actualEntry.Value, "mismatched value")
}

// run through all of the messages and make sure there aren't any panics from unimplemented methods
func (suite *ReceiverTest) TestNoPanics() {
	now := time.Now()
	for _, message := range []mqtt.SampleMessage{
		mqtt.SampleRain(now),
		mqtt.SampleTemp(now),
		mqtt.SampleSensorPause(now),
		mqtt.SampleSensorUnpause(now),
		mqtt.SampleSensorSoftReset(now),
		mqtt.SampleSensorHardReset(now),
		mqtt.SampleSensorStatus(now),
		mqtt.SampleGatewayStatus(now),
	} {
		suite.client.Publish(process(message))
	}
}

// publish a bunch of stuff to the broker
func process(msg mqtt.SampleMessage) (string, byte, bool, []byte) {
	payload, err := json.Marshal(msg.Msg)
	if err != nil {
		logrus.Error(err)
		panic("problem marshalling json")
	}
	qos := byte(viper.GetUint(configkey.MQTTQos))
	return msg.Topic, qos, false, payload
}
