// Package messenger ferries data between serial port, mqtt, and the database
package messenger

import (
	"os"
	"time"

	"github.com/ntbloom/raincounter/common/mqtt"
	"github.com/ntbloom/raincounter/sbc/database"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/raincounter/config/configkey"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Messenger receives Message from serial port, publishes to paho and stores locally
type Messenger struct {
	client paho.Client           // MQTT Client object
	db     *database.DBConnector // Database connector
	state  chan uint8            // What is the Messenger supposed to do?
	Data   chan *Message         // Actual data packets
}

// NewMessenger gets a new messenger
func NewMessenger(client paho.Client, db *database.DBConnector) *Messenger {
	state := make(chan uint8, 1)
	data := make(chan *Message, 1)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	return &Messenger{client, db, state, data}
}

// Start waits for packet to publish or to receive signal interrupt
func (m *Messenger) Start() {
	defer m.client.Disconnect(viper.GetUint(configkey.MQTTQuiescence))

	// configure status messages
	statusTimer := time.NewTicker(viper.GetDuration(configkey.MessengerStatusInterval))

	// loop until signal
	for {
		select {
		case state := <-m.state:
			switch state {
			case configkey.Kill:
				// program is exiting
				logrus.Debug("received `Closed` signal on messenger.state channel")
				statusTimer.Stop()
				return
			default:
				continue
			}
		case msg := <-m.Data:
			logrus.Debugf("received Message from serial port: %s", msg.payload)
			m.publish(msg)
		case <-statusTimer.C:
			logrus.Debug("requesting status message")
			m.sendStatus()
		}
	}
}

// Stop kills the main loop
func (m *Messenger) Stop() {
	logrus.Info("stopping messenger and closing paho connection")
	m.state <- configkey.Kill
}

// publish sends a Message over MQTT
func (m *Messenger) publish(msg *Message) {
	logrus.Tracef("sending Message over MQTT: %s", msg.payload)
	m.client.Publish(msg.topic, msg.qos, msg.retained, msg.payload)
}

// sendStatus sends a status message about the gateway and sensor at regular interval
func (m *Messenger) sendStatus() {
	// assume if this code is running that the gateway is up
	gwStatus, _ := gatewayStatusMessage()
	m.publish(gwStatus)

	sensorStatus, _ := sensorStatusMessage()
	m.publish(sensorStatus)
}

// get a status message about how the gateway is doing
func gatewayStatusMessage() (*Message, error) {
	gs := GatewayStatus{
		Topic:     mqtt.GatewayStatus,
		OK:        true,
		Timestamp: time.Now(),
	}
	msg, err := gs.Process()
	if err != nil {
		return nil, err
	}

	return &Message{
		topic:    gs.Topic,
		retained: false,
		qos:      0,
		payload:  msg,
	}, nil
}

// get a status message about how the sensor is doing
func sensorStatusMessage() (*Message, error) {
	var up bool
	port := viper.GetString(configkey.USBConnectionPort)
	_, err := os.Stat(port)
	if err != nil {
		up = false
	} else {
		up = true
	}
	ss := SensorStatus{
		Topic:     mqtt.SensorStatus,
		OK:        up,
		Timestamp: time.Now(),
	}
	msg, err := ss.Process()
	if err != nil {
		return nil, err
	}
	return &Message{
		topic:    ss.Topic,
		retained: false,
		qos:      0,
		payload:  msg,
	}, nil
}
