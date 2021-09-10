package receiver

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/raincounter/pkg/common/mqtt"
	"github.com/ntbloom/raincounter/pkg/server/webdb"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	mqttConnection paho.Client
	db             webdb.DBEntry
}

// NewReceiver creates a new Receiver struct
// The mqtt connection is created automatically and must be closed
func NewReceiver(client paho.Client) (*Receiver, error) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	db := webdb.NewPGConnector()
	recv := Receiver{
		mqttConnection: client,
		db:             db,
	}

	client.Subscribe(mqtt.RainTopic, mqtt.Qos, recv.handleRainTopic)
	client.Subscribe(mqtt.TemperatureTopic, mqtt.Qos, recv.handleTemperatureTopic)
	client.Subscribe(mqtt.GatewayStatusTopic, mqtt.Qos, recv.handleGatewayStatusMessage)
	client.Subscribe(mqtt.SensorStatusTopic, mqtt.Qos, recv.handleSensorStatusMessage)
	client.Subscribe(mqtt.SensorEventTopic, mqtt.Qos, recv.handleSensorEvent)
	return &recv, nil
}

func (r *Receiver) Close() {
	logrus.Info("disconnecting Receiver from mqtt")
	r.mqttConnection.Disconnect(viper.GetUint(configkey.MQTTQuiescence))
	logrus.Info("disconnecting Receiver from the database")
	r.db.Close()
}

func (r *Receiver) IsConnected() bool {
	return r.mqttConnection.IsConnected()
}

/* TOPIC SUBSCRIPTION CALLBACKS */

func (r *Receiver) handleGatewayStatusMessage(_ paho.Client, message paho.Message) {
	panic("not implemented!")
}

func (r *Receiver) handleSensorStatusMessage(_ paho.Client, message paho.Message) {
	panic("not implemented!")
}

func (r *Receiver) handleTemperatureTopic(_ paho.Client, message paho.Message) {
	stamp, readable, err := parseMessage(message)
	if err != nil {
		logrus.Errorf("skipping message on %s", message.Topic())
		return
	}
	temp := int(readable["TempC"].(float64))
	if err := r.db.AddTempCValue(temp, stamp); err != nil {
		logrus.Error(err)
	}
}

func (r *Receiver) handleRainTopic(_ paho.Client, message paho.Message) {
	stamp, readable, err := parseMessage(message)
	if err != nil {
		logrus.Errorf("skipping message on %s", message.Topic())
		return
	}
	mm := readable["Millimeters"].(float64)
	if err := r.db.AddRainMMEvent(mm, stamp); err != nil {
		logrus.Error(err)
	}
}

func (r *Receiver) handleSensorEvent(_ paho.Client, message paho.Message) {
	panic("not implemented!")
}

/* HELPER METHODS */

func parseMessage(msg paho.Message) (time.Time, map[string]interface{}, error) {
	var readable map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &readable); err != nil {
		logrus.Error(err)
		return time.Time{}, nil, err
	}
	stamp, err := time.Parse(configkey.TimestampFormat, readable["Timestamp"].(string))
	if err != nil {
		logrus.Error(err)
		return time.Time{}, nil, err
	}
	return stamp, readable, nil
}
