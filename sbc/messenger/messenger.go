// Package messenger ferries data between serial port, paho, and the database
package messenger

import (
	"os"
	"time"

	"github.com/ntbloom/rainbase/pkg/database"
	"github.com/ntbloom/rainbase/pkg/paho"
	"github.com/ntbloom/rainbase/pkg/timer"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Messenger receives Message from serial port, publishes to mqtt and stores locally
type Messenger struct {
	client mqtt.Client           // MQTT Client object
	db     *database.DBConnector // Database connector
	State  chan uint8            // What is the Messenger supposed to do?
	Data   chan *Message         // Actual data packets
}

// NewMessenger get a new messenger
func NewMessenger(client mqtt.Client, db *database.DBConnector) *Messenger {
	state := make(chan uint8)
	data := make(chan *Message)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	return &Messenger{client, db, state, data}
}

// Wait for packet to publish or to receive signal interrupt
func (m *Messenger) Listen() {
	defer m.client.Disconnect(viper.GetUint(configkey.MQTTQuiescence))

	// configure status messages
	statusInterval := viper.GetDuration(configkey.MessengerStatusInterval)
	statusFrequency := viper.GetDuration(configkey.MessengerStatusFrequency)
	if statusInterval > 0 {
		statusTimer := timer.NewChannelUint8Timer(statusInterval, statusFrequency, m.State, configkey.SendStatusMessage)
		go statusTimer.Loop()
	}

	// loop until signal
	for {
		select {
		case state := <-m.State:
			switch state {
			case configkey.SerialClosed:
				logrus.Debug("received `Closed` signal, closing mqtt connection")
				return
			case configkey.SendStatusMessage:
				logrus.Debug("requesting status message")
				m.SendStatus()
			}
		case msg := <-m.Data:
			logrus.Tracef("received Message from serial port: %s", msg.payload)
			m.Publish(msg)
		}
	}
}

// Publish sends a Message over MQTT
func (m *Messenger) Publish(msg *Message) {
	logrus.Tracef("sending Message over MQTT: %s", msg.payload)
	m.client.Publish(msg.topic, msg.qos, msg.retained, msg.payload)
}

// SendStatus sends a status message about the gateway and sensor at regular interval
func (m *Messenger) SendStatus() {
	// assume if this code is running that the gateway is up
	gwStatus, _ := gatewayStatusMessage()
	m.Publish(gwStatus)

	sensorStatus, _ := sensorStatusMessage()
	m.Publish(sensorStatus)
}

// get a status message about how the gateway is doing
func gatewayStatusMessage() (*Message, error) {
	gs := GatewayStatus{
		Topic:     paho.GatewayStatus,
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
		Topic:     paho.SensorStatus,
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
