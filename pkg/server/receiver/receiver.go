package receiver

import (
	"encoding/json"

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
	logrus.Error(message)
	var readable map[string]interface{}
	err := json.Unmarshal(message.Payload(), &readable)
	if err != nil {
		panic(err)
	}
	//stamp := readable["Timestamp"]
	//mm, err := strconv.ParseFloat(readable["Millimeters"], 32)

}

func (r *Receiver) handleSensorEvent() {
	panic("not implemented!")
}
