package mqtt

import (
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"
)

// Examples of sample messages for use in testing, etc.

type SampleMessage struct {
	Topic     string
	Msg       map[string]interface{}
	Timestamp time.Time
}

const (
	SampleTimestamp = "2022-09-06T21:57:32.779567444-04:00"
)

var SampleCelsius = 23

func timestamp() time.Time {
	var stamp string = string(SampleTimestamp)
	val, err := time.Parse(configkey.TimestampFormat, stamp)
	if err != nil {
		panic(err)
	}
	return val
}

func genericEventMessage(tag, value int, event string) map[string]interface{} {
	return map[string]interface{}{
		"Tag":       tag,
		"Value":     value,
		"Event":     event,
		"Timestamp": timestamp(),
	}
}

func SampleRain() SampleMessage {
	return SampleMessage{
		Topic:     RainTopic,
		Msg:       map[string]interface{}{"Millimeters": viper.GetFloat64(configkey.SensorRainMm), "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleTemp() SampleMessage {
	return SampleMessage{
		Topic:     TemperatureTopic,
		Msg:       map[string]interface{}{"TempC": SampleCelsius, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleSensorPause() SampleMessage {
	msg := genericEventMessage(tlv.Pause, tlv.PauseValue, SensorPauseEvent)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp(),
	}
}

func SampleSensorUnpause() SampleMessage {
	msg := genericEventMessage(tlv.Unpause, tlv.UnpauseValue, SensorUnpauseEvent)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp(),
	}
}

func SampleSensorSoftReset() SampleMessage {
	msg := genericEventMessage(tlv.SoftReset, tlv.SoftResetValue, SensorSoftResetEvent)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp(),
	}
}

func SampleSensorHardReset() SampleMessage {
	msg := genericEventMessage(tlv.HardReset, tlv.HardResetValue, SensorHardResetEvent)
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       msg,
		Timestamp: timestamp(),
	}
}

func SampleSensorStatus() SampleMessage {
	return SampleMessage{
		Topic:     SensorStatusTopic,
		Msg:       map[string]interface{}{"OK": true, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleGatewayStatus() SampleMessage {
	return SampleMessage{
		Topic:     GatewayStatusTopic,
		Msg:       map[string]interface{}{"OK": true, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}
