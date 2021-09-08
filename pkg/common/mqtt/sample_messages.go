package mqtt

import (
	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"
)

// Examples of sample messages for use in testing, etc.

type SampleMessage struct {
	Topic string
	Msg   map[string]interface{}
}

const (
	SampleTimestamp = "2021-09-06T21:57:32.779567444-04:00"
	SampleCelsius   = 23
)

func SampleRain() SampleMessage {
	return SampleMessage{
		Topic: RainTopic,
		Msg:   map[string]interface{}{"Millimeters": viper.GetFloat64(configkey.SensorRainMm), "Timestamp": SampleTimestamp},
	}
}

func SampleTemp() SampleMessage {
	return SampleMessage{
		Topic: TemperatureTopic,
		Msg:   map[string]interface{}{"Tempc": SampleCelsius, "Timestamp": SampleTimestamp},
	}
}

func SampleSensorPause() SampleMessage {
	return SampleMessage{
		Topic: SensorEventTopic,
		Msg:   map[string]interface{}{"Status": SensorPauseEvent, "Timestamp": SampleTimestamp},
	}
}

func SampleSensorUnpause() SampleMessage {
	return SampleMessage{
		Topic: SensorEventTopic,
		Msg:   map[string]interface{}{"Status": SensorUnpauseEvent, "Timestamp": SampleTimestamp},
	}
}

func SampleSensorSoftReset() SampleMessage {
	return SampleMessage{
		Topic: SensorEventTopic,
		Msg:   map[string]interface{}{"Status": SensorSoftResetEvent, "Timestamp": SampleTimestamp},
	}
}

func SampleSensorHardReset() SampleMessage {
	return SampleMessage{
		Topic: SensorEventTopic,
		Msg:   map[string]interface{}{"Status": SensorHardResetEvent, "Timestamp": SampleTimestamp},
	}
}

func SampleSensorStatus() SampleMessage {
	return SampleMessage{
		Topic: SensorStatusTopic,
		Msg:   map[string]interface{}{"OK": true, "Timestamp": SampleTimestamp},
	}
}

func SampleGatewayStatus() SampleMessage {
	return SampleMessage{
		Topic: GatewayStatusTopic,
		Msg:   map[string]interface{}{"OK": true, "Timestamp": SampleTimestamp},
	}
}
