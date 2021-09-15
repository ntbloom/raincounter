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
	client paho.Client
	db     webdb.DBEntry
	state  chan int
}

// NewReceiver creates a new Receiver struct
// The mqtt connection is created automatically and must be closed
func NewReceiver() (*Receiver, error) {
	client, err := mqtt.NewConnection()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	db := webdb.NewPGConnector()
	state := make(chan int)
	recv := Receiver{
		client: client,
		db:     db,
		state:  state,
	}

	qos := byte(viper.GetUint(configkey.MQTTQos))
	recv.client.Subscribe(mqtt.RainTopic, qos, recv.handleRainTopic)
	recv.client.Subscribe(mqtt.TemperatureTopic, qos, recv.handleTemperatureTopic)
	recv.client.Subscribe(mqtt.GatewayStatusTopic, qos, recv.handleGatewayStatusMessage)
	recv.client.Subscribe(mqtt.SensorStatusTopic, qos, recv.handleSensorStatusMessage)
	recv.client.Subscribe(mqtt.SensorEventTopic, qos, recv.handleSensorEvent)
	return &recv, nil
}

// Start runs the main loop, basically just waiting to be told to stop
func (r *Receiver) Start() {
	for {
		state := <-r.state
		switch state {
		case configkey.Kill:
			logrus.Debug("received `Closed` signal on receiver.state channel")
			r.Close()
			return
		default:
			logrus.Errorf("unexpected message on receiver.state channel: %d", state)
		}
	}
}

// Stop kills the main loop
func (r *Receiver) Stop() {
	logrus.Info("Stopping receiver and closing mqtt connection")
	r.state <- configkey.Kill
}

// Close closes the connection
func (r *Receiver) Close() {
	topics := []string{
		mqtt.RainTopic,
		mqtt.TemperatureTopic,
		mqtt.GatewayStatusTopic,
		mqtt.SensorStatusTopic,
		mqtt.SensorEventTopic,
	}
	r.client.Unsubscribe(topics...)
	logrus.Info("disconnecting Receiver from mqtt")
	r.client.Disconnect(viper.GetUint(configkey.MQTTQuiescence))
	logrus.Info("disconnecting Receiver from the database")
	r.db.Close()
}

func (r *Receiver) IsConnected() bool {
	return r.client.IsConnected()
}

/* TOPIC SUBSCRIPTION CALLBACKS */

func (r *Receiver) handleGatewayStatusMessage(_ paho.Client, message paho.Message) {
	go func() {
		r.processStatusMessage(message, configkey.GatewayStatus)
	}()
}

func (r *Receiver) handleSensorStatusMessage(_ paho.Client, message paho.Message) {
	go func() {
		r.processStatusMessage(message, configkey.SensorStatus)
	}()
}

func (r *Receiver) handleTemperatureTopic(_ paho.Client, message paho.Message) {
	go func() {
		stamp, readable, err := parseMessage(message)
		if err != nil {
			return
		}
		temp := int(readable["TempC"].(float64))
		if err := r.db.AddTempCValue(temp, stamp); err != nil {
			logrus.Error(err)
		}
	}()
}

func (r *Receiver) handleRainTopic(_ paho.Client, message paho.Message) {
	go func() {
		stamp, readable, err := parseMessage(message)
		if err != nil {
			return
		}
		mm := readable["Millimeters"].(float64)
		if err := r.db.AddRainMMEvent(mm, stamp); err != nil {
			logrus.Error(err)
		}
	}()
}

func (r *Receiver) handleSensorEvent(_ paho.Client, message paho.Message) {
	go func() {
		stamp, readable, err := parseMessage(message)
		if err != nil {
			return
		}
		tag := int(readable["Tag"].(float64))
		value := int(readable["Value"].(float64))
		if err := r.db.AddTagValue(tag, value, stamp); err != nil {
			return
		}
	}()
}

/* HELPER METHODS */

// send a sensor or gateway status message
func (r *Receiver) processStatusMessage(msg paho.Message, asset int) {
	stamp, _, err := parseMessage(msg)
	if err != nil {
		return
	}
	if err := r.db.AddStatusUpdate(asset, stamp); err != nil {
		logrus.Error(err)
		return
	}
}

// parse the messages and have unified error logging for all topics
func parseMessage(msg paho.Message) (time.Time, map[string]interface{}, error) {
	var readable map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &readable); err != nil {
		logrus.Errorf("skipping message on %s: %s", msg.Topic(), err)
		return time.Time{}, nil, err
	}
	stamp, err := time.Parse(configkey.TimestampFormat, readable["Timestamp"].(string))
	if err != nil {
		logrus.Errorf("skipping message on %s: %s", msg.Topic(), err)
		return time.Time{}, nil, err
	}
	return stamp, readable, nil
}
