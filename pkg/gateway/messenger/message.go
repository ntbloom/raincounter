package messenger

// Message defines what an individual message looks like

import (
	"encoding/json"
	"time"

	"github.com/ntbloom/raincounter/pkg/common/database"

	mqtt2 "github.com/ntbloom/raincounter/pkg/common/mqtt"

	configkey2 "github.com/ntbloom/raincounter/pkg/config/configkey"

	tlv2 "github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// values for static status messages
const (
	SensorPause     = "sensorPause"
	SensorUnpause   = "sensorUnpause"
	SensorSoftReset = "sensorSoftReset"
	SensorHardReset = "sensorHardReset"
)

// Payload generic type of message we'll send over MQTT
type Payload interface {
	Process() ([]byte, error)
}

// generic wrapper for all implementations of Payload
func process(p Payload) ([]byte, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	val := string(payload)
	logrus.Debug(val)
	return payload, nil
}

// SensorEvent gives static message about what's happening to the sensor
type SensorEvent struct {
	Topic     string
	Status    string
	Timestamp time.Time
}

// TemperatureEvent sends current temperature in Celsius
type TemperatureEvent struct {
	Topic     string
	TempC     int
	Timestamp time.Time
}

// RainEvent sends message about rain event
type RainEvent struct {
	Topic       string
	Millimeters string // send as a string to avoid floating point weirdness
	Timestamp   time.Time
}

// GatewayStatus sends "OK" message at regular intervals
type GatewayStatus struct {
	Topic     string
	OK        bool
	Timestamp time.Time
}

// SensorStatus sends "OK" if sensor is reachable, else "Bad"
type SensorStatus struct {
	Topic     string
	OK        bool
	Timestamp time.Time
}

// Process turn static value into mqtt payload
func (s *SensorEvent) Process() ([]byte, error) {
	return process(s)
}

// Process turn temp into mqtt payload
func (t *TemperatureEvent) Process() ([]byte, error) {
	return process(t)
}

// Process turn rain event into mqtt payload
func (r *RainEvent) Process() ([]byte, error) {
	return process(r)
}

// Process turn gateway status message into mqtt payload
func (gs *GatewayStatus) Process() ([]byte, error) {
	return process(gs)
}

// Process turn sensor status message into mqtt payload
func (ss *SensorStatus) Process() ([]byte, error) {
	return process(ss)
}

type Message struct {
	topic    string
	retained bool
	qos      byte
	payload  []byte
}

// NewMessage makes a new message from a tlv packet mqtt topic and logs the entry to the postgresql in the background
func (m *Messenger) NewMessage(packet *tlv2.TLV) (*Message, error) {
	now := time.Now()
	var event Payload
	var topic string

	switch packet.Tag {
	case tlv2.Rain:
		topic = mqtt2.RainTopic
		event = &RainEvent{
			topic,
			viper.GetString(configkey2.SensorRainMetric),
			now,
		}
		go database.MakeRainEntry(m.db)
	case tlv2.Temperature:
		topic = mqtt2.TemperatureTopic
		tempC := packet.Value
		event = &TemperatureEvent{
			topic,
			tempC,
			now,
		}
		go database.MakeTemperatureEntry(m.db, tempC)
	case tlv2.SoftReset:
		topic = mqtt2.SensorEvent
		event = &SensorEvent{
			topic,
			SensorSoftReset,
			now,
		}
		go database.MakeSoftResetEntry(m.db)
	case tlv2.HardReset:
		topic = mqtt2.SensorEvent
		event = &SensorEvent{
			topic,
			SensorHardReset,
			now,
		}
		go database.MakeHardResetEntry(m.db)
	case tlv2.Pause:
		topic = mqtt2.SensorEvent
		event = &SensorEvent{
			topic,
			SensorPause,
			now,
		}
		go database.MakePauseEntry(m.db)
	case tlv2.Unpause:
		topic = mqtt2.SensorEvent
		event = &SensorEvent{
			topic,
			SensorUnpause,
			now,
		}
		go database.MakeUnpauseEntry(m.db)
	default:
		logrus.Errorf("unsupported tag %d", packet.Tag)
		return nil, nil
	}

	payload, err := event.Process()
	if err != nil {
		return nil, err
	}
	msg := Message{
		topic:    topic,
		retained: false,
		qos:      byte(viper.GetInt(configkey2.MQTTQos)),
		payload:  payload,
	}
	return &msg, nil
}
