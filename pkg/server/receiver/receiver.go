package receiver

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/raincounter/pkg/common/database"
	"github.com/ntbloom/raincounter/pkg/common/mqtt"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	mqttConnection  paho.Client
	sqliteConection *database.Sqlite
}

// NewReceiver creates a new Receiver struct
// mqtt connection is created automatically
func NewReceiver(mqttConfig *mqtt.BrokerConfig, databasePath string, clobber bool) (*Receiver, error) {
	s, err := database.NewSqlite(databasePath, clobber)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	m, err := mqtt.NewConnection(mqttConfig)
	if token := m.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &Receiver{
		mqttConnection:  m,
		sqliteConection: s,
	}, nil
}

func (r *Receiver) IsConnected() bool {
	return r.mqttConnection.IsConnected()
}

func (r *Receiver) receiveGatewayStatusMessage() {
	logrus.Error("not implemented!")
}
