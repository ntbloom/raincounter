package mqtt

import (
	"time"

	"github.com/ntbloom/raincounter/pkg/rainbase/tlv"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"
)

// Examples of sample messages for use in testing, etc.

// SampleMessage is a dummy message
type SampleMessage struct {
	Topic     string
	Msg       map[string]interface{}
	Timestamp time.Time
}

// SampleCelsius is a random temperature value picked for no reason
var SampleCelsius = 23

func genericEventMessage(tag, value int, event string, timestamp time.Time) map[string]interface{} {
	return map[string]interface{}{
		"Tag":       tag,
		"Value":     value,
		"Event":     event,
		"Timestamp": timestamp,
	}
}

// SampleRain is a test mqtt message for a rain event
func SampleRain(timestamp time.Time) SampleMessage {
	return SampleMessage{
		Topic:     RainTopic,
		Msg:       map[string]interface{}{"Millimeters": viper.GetFloat64(configkey.SensorRainMm), "Timestamp": timestamp},
		Timestamp: timestamp,
	}
}

// SampleTemp is a test mqtt message for a temperature measurement in C
func SampleTemp(timestamp time.Time) SampleMessage {
	return SampleMessage{
		Topic:     TemperatureTopic,
		Msg:       map[string]interface{}{"TempC": SampleCelsius, "Timestamp": timestamp},
		Timestamp: timestamp,
	}
}

// SampleSensorPause is a test mqtt message for a pause event
func SampleSensorPause(timestamp time.Time) SampleMessage {
	msg := genericEventMessage(tlv.Pause, tlv.PauseValue, SensorPauseEvent, timestamp)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp,
	}
}

// SampleSensorUnpause is a test mqtt message for an unpause event
func SampleSensorUnpause(timestamp time.Time) SampleMessage {
	msg := genericEventMessage(tlv.Unpause, tlv.UnpauseValue, SensorUnpauseEvent, timestamp)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp,
	}
}

// SampleSensorSoftReset is a test mqtt message for a soft reset event
func SampleSensorSoftReset(timestamp time.Time) SampleMessage {
	msg := genericEventMessage(tlv.SoftReset, tlv.SoftResetValue, SensorSoftResetEvent, timestamp)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp,
	}
}

// SampleSensorHardReset is a test mqtt message for a hard reset event
func SampleSensorHardReset(timestamp time.Time) SampleMessage {
	msg := genericEventMessage(tlv.HardReset, tlv.HardResetValue, SensorHardResetEvent, timestamp)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp,
	}
}

// SampleSensorStatus is a test mqtt message for a sensor status message
func SampleSensorStatus(timestamp time.Time) SampleMessage {
	return SampleMessage{
		Topic:     SensorStatusTopic,
		Msg:       map[string]interface{}{"OK": true, "Timestamp": timestamp},
		Timestamp: timestamp,
	}
}

// SampleGatewayStatus is a test mqtt message for a gateway status message
func SampleGatewayStatus(timestamp time.Time) SampleMessage {
	return SampleMessage{
		Topic:     GatewayStatusTopic,
		Msg:       map[string]interface{}{"OK": true, "Timestamp": timestamp},
		Timestamp: timestamp,
	}
}
