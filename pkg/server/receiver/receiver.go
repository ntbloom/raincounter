package receiver

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/raincounter/pkg/gateway/localdb"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	mqttConnection  paho.Client
	sqliteConection *localdb.Sqlite
}

// NewReceiver creates a new Receiver struct
// mqtt connection is created automatically
func NewReceiver(client paho.Client, databasePath string, clobber bool) (*Receiver, error) {
	s, err := localdb.NewSqlite(databasePath, clobber)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}

	return &Receiver{
		mqttConnection:  client,
		sqliteConection: s,
	}, nil
}

func (r *Receiver) IsConnected() bool {
	return r.mqttConnection.IsConnected()
}

func (r *Receiver) handleGatewayStatusMessage() {
	logrus.Error("not implemented!")
}

func (r *Receiver) handleSensorStatusMessage() {
	logrus.Error("not implemented!")
}

func (r *Receiver) handleTemperatureMessage() {
	logrus.Error("not implemented!")
}

func (r *Receiver) handleRainTopic() {
	logrus.Error("not implemented!")
}

func (r *Receiver) handleSensorEvent() {
	logrus.Error("not implemented!")
}
