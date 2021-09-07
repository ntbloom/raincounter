package mqtt

// Examples of sample messages for use in testing, etc.

type SampleMessage struct {
	Topic string
	Msg   map[string]interface{}
}

const (
	SampleTimestamp = "2021-09-06T21:57:32.779567444-04:00"
	SampleCelsius   = 23
)

var SampleRain = SampleMessage{
	Topic: RainTopic,
	Msg:   map[string]interface{}{"Millimeters": "0.2794", "Timestamp": SampleTimestamp},
}

var SampleTemp = SampleMessage{
	Topic: TemperatureTopic,
	Msg:   map[string]interface{}{"Tempc": SampleCelsius, "Timestamp": SampleTimestamp},
}

var SampleSensorPause = SampleMessage{
	Topic: SensorEventTopic,
	Msg:   map[string]interface{}{"Status": SensorPauseEvent, "Timestamp": SampleTimestamp},
}

var SampleSensorUnpause = SampleMessage{
	Topic: SensorEventTopic,
	Msg:   map[string]interface{}{"Status": SensorUnpauseEvent, "Timestamp": SampleTimestamp},
}

var SampleSensorSoftReset = SampleMessage{
	Topic: SensorEventTopic,
	Msg:   map[string]interface{}{"Status": SensorSoftResetEvent, "Timestamp": SampleTimestamp},
}

var SampleSensorHardReset = SampleMessage{
	Topic: SensorEventTopic,
	Msg:   map[string]interface{}{"Status": SensorHardResetEvent, "Timestamp": SampleTimestamp},
}

var SampleSensorStatus = SampleMessage{
	Topic: SensorStatusTopic,
	Msg:   map[string]interface{}{"OK": true, "Timestamp": SampleTimestamp},
}

var SampleGatewayStatus = SampleMessage{
	Topic: GatewayStatusTopic,
	Msg:   map[string]interface{}{"OK": true, "Timestamp": SampleTimestamp},
}
