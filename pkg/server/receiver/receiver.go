package receiver

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/raincounter/pkg/gateway/localdb"
	"github.com/ntbloom/raincounter/pkg/server/webdb"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	mqttConnection  paho.Client
	sqliteConection *localdb.LocalDB
	db              webdb.DBEntry
}

// NewReceiver creates a new Receiver struct
// mqtt connection is created automatically
func NewReceiver(client paho.Client, databasePath string, clobber bool) (*Receiver, error) {
	s, err := localdb.NewLocalDB(databasePath, clobber)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}

	var db webdb.DBEntry
	db, err = webdb.NewWebSqlite(databasePath, true)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &Receiver{
		mqttConnection:  client,
		sqliteConection: s,
		db:              db,
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
