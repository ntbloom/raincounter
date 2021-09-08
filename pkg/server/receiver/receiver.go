package receiver

import (
	"encoding/json"
	"time"

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
// mqtt connection is created automatically
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
	return &recv, nil
}

func (r *Receiver) IsConnected() bool {
	return r.mqttConnection.IsConnected()
}

func (r *Receiver) handleGatewayStatusMessage() {
	panic("not implemented!")
}

func (r *Receiver) handleSensorStatusMessage() {
	panic("not implemented!")
}

func (r *Receiver) handleTemperatureMessage() {
	panic("not implemented!")
}

func (r *Receiver) handleRainTopic(client paho.Client, message paho.Message) {
	var readable map[string]interface{}
	if err := json.Unmarshal(message.Payload(), &readable); err != nil {
		logrus.Error(err)
		return
	}
	stamp, err := time.Parse(configkey.TimestampFormat, readable["Timestamp"].(string))
	if err != nil {
		logrus.Error(err)
		return
	}
	mm := readable["Millimeters"].(float64)
	if err := r.db.AddRainMMEvent(mm, stamp); err != nil {
		logrus.Error(err)
		return
	}
}

func (r *Receiver) handleSensorEvent() {
	panic("not implemented!")
}
