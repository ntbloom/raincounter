package receiver

import (
	paho "github.com/eclipse/paho.mqtt.golang"
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

	return &Receiver{
		mqttConnection: client,
		db:             db,
	}, nil
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

func (r *Receiver) handleRainTopic() {
	panic("not implemented!")
}

func (r *Receiver) handleSensorEvent() {
	panic("not implemented!")
}
