package messenger

// Message defines what an individual message looks like

import (
	"encoding/json"
	"time"

	"github.com/ntbloom/raincounter/pkg/common/database"
	"github.com/ntbloom/raincounter/pkg/common/mqtt"
	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/ntbloom/raincounter/pkg/gateway/tlv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	//val := string(payload)
	//logrus.Tracef("processing payload %s: ", val)
	return payload, nil
}

// SensorEvent gives static message about what's happening to the sensor
type SensorEvent struct {
	Status    string
	Timestamp time.Time
}

// TemperatureEvent sends current temperature in Celsius
type TemperatureEvent struct {
	TempC     int
	Timestamp time.Time
}

// RainEvent sends message about rain event
type RainEvent struct {
	Millimeters string // send as a string to avoid floating point weirdness
	Timestamp   time.Time
}

// GatewayStatus sends "OK" message at regular intervals
type GatewayStatus struct {
	OK        bool
	Timestamp time.Time
}

// SensorStatus sends "OK" if sensor is reachable, else "Bad"
type SensorStatus struct {
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
func (m *Messenger) NewMessage(packet *tlv.TLV) (*Message, error) {
	now := time.Now()
	var event Payload
	var topic string

	switch packet.Tag {
	case tlv.Rain:
		topic = mqtt.RainTopic
		event = &RainEvent{
			viper.GetString(configkey.SensorRainMm),
			now,
		}
		go database.MakeRainTallyEntry(m.db)
	case tlv.Temperature:
		topic = mqtt.TemperatureTopic
		tempC := packet.Value
		event = &TemperatureEvent{
			tempC,
			now,
		}
		go database.MakeTemperatureEntry(m.db, tempC)
	case tlv.SoftReset:
		topic = mqtt.SensorEventTopic
		event = &SensorEvent{
			mqtt.SensorSoftResetEvent,
			now,
		}
		go database.MakeSoftResetEntry(m.db)
	case tlv.HardReset:
		topic = mqtt.SensorEventTopic
		event = &SensorEvent{
			mqtt.SensorHardResetEvent,
			now,
		}
		go database.MakeHardResetEntry(m.db)
	case tlv.Pause:
		topic = mqtt.SensorEventTopic
		event = &SensorEvent{
			mqtt.SensorPauseEvent,
			now,
		}
		go database.MakePauseEntry(m.db)
	case tlv.Unpause:
		topic = mqtt.SensorEventTopic
		event = &SensorEvent{
			mqtt.SensorUnpauseEvent,
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
		qos:      byte(viper.GetInt(configkey.MQTTQos)),
		payload:  payload,
	}
	logrus.Tracef("sending message, topic=%s, payload=%s", topic, payload)
	return &msg, nil
}
