package mqtt

import (
	"time"

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
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       map[string]interface{}{"Status": SensorPauseEvent, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleSensorUnpause() SampleMessage {
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       map[string]interface{}{"Status": SensorUnpauseEvent, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleSensorSoftReset() SampleMessage {
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       map[string]interface{}{"Status": SensorSoftResetEvent, "Timestamp": timestamp()},
		Timestamp: timestamp(),
	}
}

func SampleSensorHardReset() SampleMessage {
	return SampleMessage{
		Topic:     SensorEventTopic,
		Msg:       map[string]interface{}{"Status": SensorHardResetEvent, "Timestamp": timestamp()},
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
